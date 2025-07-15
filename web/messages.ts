import * as banterpb from '/m/banter/pb/banter_pb.js';
import { Cfg } from './controller.js';
import { UpdatingControlPanel } from '/tk.js';
import { BanterPlayground, help } from './banter_playground.js';

class BanterMessages extends UpdatingControlPanel<banterpb.Config> {
    private _table: HTMLTableElement;
    private _button: HTMLButtonElement;
    private _editor: BanterPlayground;

    private _rows: BanterRow[] = [];

    constructor(cfg: Cfg) {
        super({ title: 'Messages', help, data: cfg });

        this.innerHTML = `
<table></table>
<button> + </button>
`;
        this._table = this.querySelector('table');
        this._button = this.querySelector('button');
        this._button.addEventListener('click', () => this._handleNew());
        this._editor = new BanterPlayground(cfg);
        this.appendChild(this._editor);
    }

    update(cfg: banterpb.Config) {
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
            row.onedit = () => this._editor.edit(banter);
            row.onsave = (banter) => {
                let cfg = this.last.clone();
                cfg.banters[i] = banter;
                this._save(cfg);
            };
            row.ondelete = () => this._delete(i);
            return row;
        });
        this._rows.forEach((row) => {
            this._table.appendChild(row);
        })
        this._editor.close();
    }

    private _handleNew() {
        this._editor.edit(new banterpb.Banter());
    }

    private _delete(i: number) {
        let cfg = this.last
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
        this.save(cfg);
    }
}
customElements.define('banter-messages', BanterMessages, { extends: 'fieldset' });

class BanterRow extends HTMLTableRowElement {
    private _input_command: HTMLInputElement;
    private _input_text: HTMLInputElement;
    private _check_enabled: HTMLInputElement;
    private _check_random: HTMLInputElement;
    private _button_edit: HTMLButtonElement;
    private _button_delete: HTMLButtonElement;
    private _orig: banterpb.Banter;
    onedit: () => void = () => { };
    onsave: (banter: banterpb.Banter) => void = () => { };
    ondelete: () => void = () => { };

    constructor() {
        super();
        this.innerHTML = `
<td><input id="command" type="text" size="16" disabled /></td>
<td><input id="text" type="text" size="48" disabled /></td>
<td><input id="enabled" type="checkbox" /></td>
<td><input id="random" type="checkbox" /></td>
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
        this._button_edit.onclick = () => this.onedit();
        this._button_delete = this.querySelector('#btn-delete');
        this._button_delete.onclick = () => this.ondelete();

        this._check_enabled.addEventListener('change', () => this.save());
        this._check_random.addEventListener('change', () => this.save());
    }

    set banter(banter: banterpb.Banter) {
        this._orig = banter.clone();
        this._input_command.value = banter.command;
        this._input_text.value = banter.text;
        this._check_enabled.checked = !banter.disabled;
        this._check_random.checked = banter.random;

        if (banter.text.includes('{{')) {
            this._check_random.disabled = true;
        }
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