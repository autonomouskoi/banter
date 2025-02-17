import * as banterpb from '/m/banter/pb/banter_pb.js';
import { GloballyStyledHTMLElement } from "/global-styles.js";
import { Cfg } from './controller.js';

class BanterMessages extends GloballyStyledHTMLElement {
    private _table: HTMLTableElement;
    private _button: HTMLButtonElement;

    private _rows: BanterRow[] = [];
    private _cfg: Cfg;

    constructor(cfg: Cfg) {
        super();
        this._cfg = cfg;

        this.shadowRoot.innerHTML = `
<fieldset>
<legend>Messages</legend>

<table></table>
<button> + </button>
</fieldset>
`;
        this._table = this.shadowRoot.querySelector('table');
        this._button = this.shadowRoot.querySelector('button');
        this._button.addEventListener('click', () => this._handleNew());

        this._cfg.subscribe((newCfg) => this.config = newCfg);
        this.config = this._cfg.last;
    }

    set config(cfg: banterpb.Config) {
        this._table.innerHTML = `
<tr>
    <th>Command</th>
    <th>Text</th>
    <th>Enabled</th>
    <th>Random</th>
    <th>Edit</th>
</tr>
`;
        this._rows = cfg.banters.map((banter, i) => {
            let row = new BanterRow();
            row.banter = banter;
            row.oncancel = () => this._cancelEditing();
            row.onedit = () => this._setEditing(i);
            row.onsave = (banter) => {
                let cfg = this._cfg.last.clone();
                cfg.banters[i] = banter;
                this._save(cfg);
            };
            row.ondelete = () => this._delete(i);
            return row;
        });
        this._rows.forEach((row) => {
            this._table.appendChild(row);
        })
    }

    private _handleNew() {
        let row = new BanterRow();
        row.banter = new banterpb.Banter();
        row.oncancel = () => {
            this._cancelEditing();
            this.config = this._cfg.last;
        };
        row.onsave = (banter) => {
            let cfg = this._cfg.last.clone(); 
            cfg.banters.push(banter);
            this._save(cfg);
        };
        this._rows.push(row);
        this._table.appendChild(row);
        row.startEdit();
        this._setEditing(this._rows.length - 1);
    }

    private _setEditing(idx: number) {
        this._button.disabled = true;
        this._rows.forEach((row, i) => {
            row.disabled = i != idx;
        })
    }

    private _cancelEditing() {
        this._button.disabled = false;
        this._rows.forEach((row) => {
            row.disabled = false;
        });
    }

    private _delete(i: number) {
        let cfg = this._cfg.last
        if (i < 0 || i >= cfg.banters?.length) {
            return;
        }
        let banter = cfg.banters[i];
        if (!confirm(`Delete ${banter.command}?`)) {
            return;
        }
        cfg = cfg.clone();
        cfg.banters.splice(i, 1);
        this._save(cfg);
    }

    private _save(cfg: banterpb.Config) {
        if (this._rows.find((row) => !row.isValid)) {
            return;
        }
        this._cfg.save(cfg);
    }
}
customElements.define('banter-messages', BanterMessages);

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

export { BanterMessages };