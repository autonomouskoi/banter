import { bus, enumName } from "/bus.js";
import * as buspb from "/pb/bus/bus_pb.js";
import * as banterpb from '/m/banter/pb/banter_pb.js';

const TOPIC_BANTER_COMMAND = enumName(banterpb.BusTopic, banterpb.BusTopic.BANTER_COMMAND);


class EventRow extends HTMLTableRowElement {
    private _input_text: HTMLInputElement;
    private _check_enabled: HTMLInputElement;
    private _button_edit: HTMLButtonElement;
    private _button_cancel: HTMLButtonElement;
    private _orig: banterpb.EventSettings = new banterpb.EventSettings({ enabled: false, text: '' });
    private _key: string;
    oncancel: () => void = () => { };
    onedit: () => void = () => { };
    onsave: (settings: banterpb.EventSettings) => void = () => { };

    constructor({ name, settings }: { name: string, settings: banterpb.EventSettings }) {
        super();
        this.innerHTML = `
<td>${name}</td>
<td><input id="enabled" type="checkbox" disabled /></td>
<td><input id="text" type="text" size="48" disabled /></td>
<td><button id="btn-edit">Edit</button></td>
<td><button id="btn-cancel" disabled >Cancel</button></td>
`;
        this._orig = settings ?? new banterpb.EventSettings();
        this._key = name;
        this._input_text = this.querySelector('#text');
        this._input_text.value = this._orig.text;
        this._check_enabled = this.querySelector('#enabled');
        this._check_enabled.checked = this._orig.enabled;
        this._button_edit = this.querySelector('#btn-edit');
        this._button_edit.onclick = () => this.startEdit();
        this._button_cancel = this.querySelector('#btn-cancel');
        this._button_cancel.onclick = () => this.cancelEdit();
    }

    get key(): string {
        return this._key;
    }

    private set editable(enabled: boolean) {
        this._input_text.disabled = !enabled;
        this._check_enabled.disabled = !enabled;
        this._button_cancel.disabled = !enabled;
    }

    set disabled(disabled: boolean) {
        this._button_edit.disabled = disabled;
        this._check_enabled.disabled = disabled;
        this._input_text.disabled = disabled;
        this._button_cancel.disabled = disabled;
    }

    startEdit() {
        this.editable = true;

        this._button_edit.innerText = 'Save';
        this._button_edit.onclick = () => this.save();
        this.onedit();
    }

    cancelEdit() {
        this._check_enabled.checked = this._orig.enabled;
        this._input_text.value = this._orig.text;
        this.editable = false;
        this._button_edit.innerText = 'Edit';
        this._button_edit.onclick = () => this.startEdit();
        this.oncancel();
    }

    save() {
        let newSettings = this._orig.clone();
        newSettings.enabled = this._check_enabled.checked;
        newSettings.text = this._input_text.value;
        this.onsave(newSettings);
    }
}
customElements.define('banter-event-row', EventRow, { extends: 'tr' });

class BanterRow extends HTMLTableRowElement {
    private _input_command: HTMLInputElement;
    private _input_text: HTMLInputElement;
    private _check_enabled: HTMLInputElement;
    private _check_random: HTMLInputElement;
    private _button_edit: HTMLButtonElement;
    private _button_delete: HTMLButtonElement;
    private _orig: banterpb.Banter;
    oncancel: () => void = () => { };
    onedit: () => void = () => { };
    onsave: (banter: banterpb.Banter) => void = () => { };
    ondelete: () => void = () => { };

    constructor() {
        super();
        this.innerHTML = `
<td><input id="command" type="text" size="16" disabled /></td>
<td><input id="text" type="text" size="48" disabled /></td>
<td><input id="enabled" type="checkbox" disabled /></td>
<td><input id="random" type="checkbox" disabled /></td>
<td>
    <button id="btn-edit">Edit</button>
    <button id="btn-delete">Delete</button>
</td>
`;

        this._input_command = this.querySelector('#command');
        this._input_text = this.querySelector('#text');
        this._check_enabled = this.querySelector('#enabled');
        this._check_random = this.querySelector('#random');
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
        this._check_random.checked = banter.random;
    }

    private set editable(enabled: boolean) {
        this._input_command.disabled = !enabled;
        this._input_text.disabled = !enabled;
        this._check_random.disabled = !enabled;
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
        newBanter.random = this._check_random.checked;
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
    private _table_events: HTMLTableElement;
    private _button_new: HTMLButtonElement;
    private _banter_rows: BanterRow[] = [];
    private _event_rows: { [key: string]: EventRow} = {};

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
<legend>Random Command Configuration</legend>

<label for="input-interval-seconds">Random Command Interval (seconds)</label>
<input id="input-interval-seconds" type="text"
    inputmode="numeric" pattern="\d+" size="4" placeholder="300"
    title="How long between random chat commands""
/>

<label for="input-cooldown-seconds">Random Command Cooldown (seconds)</label>
<input id="input-cooldown-seconds" type="type"
    inputmode="numeric" pattern="\d+" size="4" placeholder="900"
    title="Minimum time before repeating a given random command"
/>

</fieldset>
<fieldset>
<legend>Messages</legend>
    <table id="banters"></table>
    <button id="btn-new-banter"> + </button>
</fieldset>

<fieldset>
<legend>Events</legend>
    <table id="events"></table>
</fieldset>
`;
        this._input_cooldown = this.shadowRoot.querySelector('#input-cooldown-seconds');
        this._input_interval = this.shadowRoot.querySelector('#input-interval-seconds');
        let updateCfg = () => { this._save() };
        this._input_cooldown.onchange = updateCfg;
        this._input_interval.onchange = updateCfg;
        this._table_banters = this.shadowRoot.querySelector('#banters');
        this._table_events = this.shadowRoot.querySelector('#events');
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
    <th>Random</th>
    <th>Edit</th>
</tr>
`;
        this._banter_rows = this._config.banters.map((banter, i) => {
            let row = new BanterRow();
            row.banter = banter;
            row.oncancel = () => this._cancelEditingBanters();
            row.onedit = () => this._setEditingBanter(i);
            row.onsave = (banter) => {
                this._config.banters[i] = banter;
                this._save();
            };
            row.ondelete = () => this._deleteBanter(i);
            return row;
        });
        this._banter_rows.forEach((row: BanterRow) => {
            this._table_banters.appendChild(row);
        });
        this._table_events.innerHTML = `
<tr>
    <th>Event</th>
    <th>Enabled</th>
    <th>Text</th>
    <th>Edit</th>
</tr>
`;
        this._event_rows = {};
        let raid = new EventRow({ name: 'Raid', settings: this._config.channelRaid });
        raid.onsave = (es) => {
            this._config.channelRaid = es;
            this._save();
        }
        this._event_rows['Raid'] = raid;
        let follow = new EventRow({ name: 'Follow', settings: this._config.channelFollow });
        follow.onsave = (es) => {
            this._config.channelFollow = es;
            this._save();
        }
        this._event_rows['Follow'] = follow;
        let cheer  = new EventRow({ name: 'Cheer', settings: this._config.channelCheer });
        cheer.onsave = (es) => {
            this._config.channelCheer = es;
            this._save();
        }
        this._event_rows['Cheer'] = cheer;
        let redeem = new EventRow({ name: 'Custom Points Redeem', settings: this._config.channelPointsCustomRedeem });
        redeem.onsave = (es) => {
            this._config.channelPointsCustomRedeem = es;
            this._save();
        }
        this._event_rows['CustomRedeem'] = redeem;

        Object.keys(this._event_rows).forEach((key) => {
            let er = this._event_rows[key];
            er.oncancel = () => this._cancelEditingEvents();
            er.onedit = () => this._setEditingEvent(key);
            this._table_events.appendChild(er);
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
            this._cancelEditingBanters();
            this._config.banters.pop();
            this._update();
        }
        banterRow.onsave = (banter) => {
            this._config.banters[this._config.banters.length - 1] = banter;
            this._save();
        };
        this._table_banters.appendChild(banterRow);
        banterRow.startEdit();

        this._setEditingBanter(this._config.banters.length - 1);
    }

    private _setEditingBanter(idx: number) {
        this._button_new.disabled = true;
        this._banter_rows.forEach((row, i) => {
            row.disabled = i != idx;
        });
    }

    private _setEditingEvent(editKey: string) {
        Object.keys(this._event_rows).forEach((key) => {
            let er = this._event_rows[key];
            er.disabled = editKey !== key;
        });
    }

    private _cancelEditingBanters() {
        this._button_new.disabled = false;
        this._banter_rows.forEach((row, i) => {
            row.disabled = false;
        });
    }

    private _cancelEditingEvents() {
        Object.keys(this._event_rows).forEach((key) => {
            this._event_rows[key].disabled = false;
        });
    }

    private _deleteBanter(i: number) {
        if (i < 0 || i >= this._config.banters.length) {
            return;
        }
        let banter = this._config.banters[i];
        if (!confirm(`Delete ${banter.command}?`)) {
            return;
        }
        this._config.banters.splice(i, 1);
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
            this._cancelEditingBanters();
            this._cancelEditingEvents();
            let csResp = banterpb.ConfigSetResponse.fromBinary(reply.message);
            this.config = csResp.config;
        });
    }
}
customElements.define('banter-cfg', Config);


export { Config };