import { bus, enumName } from "/bus.js";
import * as buspb from "/pb/bus/bus_pb.js";
import * as banterpb from '/m/banter/pb/banter_pb.js';

const TOPIC_BANTER_COMMAND = enumName(banterpb.BusTopic, banterpb.BusTopic.BANTER_ANNOUNCE_COMMAND);

class BanterRow extends HTMLTableRowElement {
    private _input_command: HTMLInputElement;
    private _input_text: HTMLInputElement;
    private _check_enabled: HTMLInputElement;
    private _check_announce: HTMLInputElement;
    private _button_edit: HTMLButtonElement;
    private _button_delete: HTMLButtonElement;
    private _orig: banterpb.Banter;
    oncancel: () => void = () => {};
    onedit: () => void = () => {};
    onsave: (banter: banterpb.Banter) => void = () => {};
    ondelete: () => void = () => {};

    constructor() {
        super();
        this.innerHTML = `
<td><input id="command" type="text" size="16" disabled /></td>
<td><input id="text" type="text" size="48" disabled /></td>
<td><input id="enabled" type="checkbox" disabled /></td>
<td><input id="announce" type="checkbox" disabled /></td>
<td>
    <button id="btn-edit">Edit</button>
    <button id="btn-delete">Delete</button>
</td>
`;

        this._input_command = this.querySelector('#command');
        this._input_text = this.querySelector('#text');
        this._check_enabled = this.querySelector('#enabled');
        this._check_announce = this.querySelector('#announce');
        this._button_edit = this.querySelector('#btn-edit');
        this._button_edit.onclick = () => this.startEdit();
        this._button_delete = this.querySelector('#btn-delete');
        this._button_delete.onclick = () => this.ondelete();
    }

    set banter(banter: banterpb.Banter) {
        this._orig = banter.clone();
        this._input_command.value = banter.command;
        this._input_text.value = banter.text;
        this._check_enabled.checked = !banter.disabled;
        this._check_announce.checked = banter.announce;
    }

    private set editable(enabled: boolean) {
        this._input_command.disabled = !enabled;
        this._input_text.disabled = !enabled;
        this._check_announce.disabled = !enabled;
        this._check_enabled.disabled = !enabled;
    }

    set disabled(disabled: boolean) {
        this._button_edit.disabled = disabled;
        this._button_delete.disabled = disabled;
    }

    startEdit() {
        this.editable = true;

        this._button_edit.innerText = 'Save';
        this._button_edit.onclick = () => this.save();
        this._button_delete.innerText = 'Cancel';
        this._button_delete.onclick = () => this.cancelEdit();
        this.onedit();
    }

    cancelEdit() {
        this.banter = this._orig;
        this.editable = false;
        this._button_edit.innerText = 'Edit';
        this._button_edit.onclick = () => this.startEdit();
        this._button_delete.innerText = 'Delete';
        this._button_delete.onclick = () => this.ondelete();
        this.oncancel();
    }

    save() {
        let newBanter = this._orig.clone();
        newBanter.command = this._input_command.value;
        newBanter.text = this._input_text.value;
        newBanter.announce = this._check_announce.checked;
        newBanter.disabled = !this._check_enabled.checked;
        this.onsave(newBanter);
    }

    isValid(): boolean {
        return true;
    }
}
customElements.define('banter-row', BanterRow, { extends: 'tr' });

class Config extends HTMLElement {
    private _input_interval: HTMLInputElement;
    private _input_cooldown: HTMLInputElement;
    private _config: banterpb.Config;
    private _table_banters: HTMLTableElement;
    private _button_new: HTMLButtonElement;
    private _banter_rows: BanterRow[] = [];

    constructor() {
        super();
        this.attachShadow({ mode: 'open' });
        this.shadowRoot.innerHTML = `
<style>
#fs-announce-cfg {
    display: grid;
    grid-template-columns: max-content max-content;
    column-gap: 1rem;
}
</style>
<fieldset id="fs-announce-cfg">
<legend>Announcement Configuration</legend>

<label for="input-interval-seconds">Announce Interval (seconds)</label>
<input id="input-interval-seconds" type="text"
    inputmode="numeric" pattern="\d+" size="4" placeholder="300"
    title="How long between chat announcments""
/>

<label for="input-cooldown-seconds">Announce Cooldown (seconds)</label>
<input id="input-cooldown-seconds" type="type"
    inputmode="numeric" pattern="\d+" size="4" placeholder="900"
    title="Minimum time before repeating a given announcement"
/>

</fieldset>
<fieldset>
<legend>Messages</legend>
    <table id="banters"></table>
    <button id="btn-new-banter"> + </button>
</fieldset>
`;
        this._input_cooldown = this.shadowRoot.querySelector('#input-cooldown-seconds');
        this._input_interval = this.shadowRoot.querySelector('#input-interval-seconds');
        let updateCfg = () => { this._save() };
        this._input_cooldown.onchange = updateCfg;
        this._input_interval.onchange = updateCfg;
        this._table_banters = this.shadowRoot.querySelector('#banters');
        this._button_new = this.shadowRoot.querySelector('#btn-new-banter');
        this._button_new.onclick = () => this.handleNew();
    }

    private _update() {
        this._input_cooldown.value = this._config.cooldownSeconds.toString();
        this._input_interval.value = this._config.intervalSeconds.toString();

        this._table_banters.innerHTML = `
<tr>
    <th>Command</th>
    <th>Text</th>
    <th>Enabled</th>
    <th>Announced</th>
    <th>Edit</th>
</tr>
`;
        this._banter_rows = this._config.banters.map((banter, i) => {
            let row = new BanterRow();
            row.banter = banter;
            row.oncancel = () => this._cancelEditing();
            row.onedit = () => this._setEditing(i);
            row.onsave = (banter) => {
                this._config.banters[i] = banter;
                this._save();
            };
            row.ondelete = () => this._delete(i);
            return row;
        });
        this._banter_rows.forEach((row: BanterRow) => {
            this._table_banters.appendChild(row);
        });

    }

    set config(config: banterpb.Config) {
        this._config = config;
        this._update();
    }

    handleNew() {
        let newBanter = new banterpb.Banter();
        this._config.banters.push(newBanter);

        // set callbacks
        let banterRow = new BanterRow();
        banterRow.banter = newBanter;
        banterRow.oncancel = () => {
            this._cancelEditing();
            this._config.banters.pop();
            this._update();
        }
        banterRow.onsave = (banter) => {
            this._config.banters[this._config.banters.length - 1] = banter;
            this._save();
        };
        this._table_banters.appendChild(banterRow);
        banterRow.startEdit();

        this._setEditing(this._config.banters.length - 1);
    }

    private _setEditing(idx: number) {
        this._button_new.disabled = true;
        this._banter_rows.forEach((row, i) => {
            row.disabled = i != idx;
        });
    }

    private _cancelEditing() {
        this._button_new.disabled = false;
        this._banter_rows.forEach((row, i) => {
            row.disabled = false;
        });
    }
    
    private _delete(i: number) {
        if (i < 0 || i >= this._config.banters.length) {
            console.log(`DEBUG delete ${i}/${this._config.banters.length}`);
            return;
        }
        let banter = this._config.banters[i];
        if (!confirm(`Delete ${banter.command}?`)) {
            return;
        }
        this._config.banters.splice(i,1);
        this._save();
    }

    private _save() {
        if (this._banter_rows.find((row) => !row.isValid)) {
            return;
        }
        let msg = new buspb.BusMessage();
        msg.topic = TOPIC_BANTER_COMMAND;
        msg.type = banterpb.MessageTypeCommand.CONFIG_SET_REQ;
        let csReq = new banterpb.ConfigSetRequest(); 
        csReq.config = this._config;
        csReq.config.intervalSeconds = parseInt(this._input_interval.value);
        csReq.config.cooldownSeconds = parseInt(this._input_cooldown.value);
        msg.message = csReq.toBinary();
        bus.sendWithReply(msg, (reply) => {
            if (reply.error) {
                alert(reply.error.detail);
                return;
            }
            this._cancelEditing();
            let csResp = banterpb.ConfigSetResponse.fromBinary(reply.message);
            this.config = csResp.config;
        });
    }
}
customElements.define('banter-cfg', Config);


export { Config };