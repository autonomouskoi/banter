import * as banterpb from '/m/banter/pb/banter_pb.js';
import { Cfg } from './controller.js';
import { UpdatingControlPanel } from '/tk.js';

let help = document.createElement('div');
help.innerHTML = `
<p>
Banter can send custom messages to chat when certain events occur. The messages can contain
special placeholder values that will incorporate information about the event into the chat
message.
</p>

<h3>Raid Placeholders</h3>
<dl>
    <dt><code>{{ .FromBroadcaster.Login }}</code></dt>
        <dd>The login name of the raider</dd>
    <dt><code>{{ .FromBroadcaster.Name }}</code></dt>
        <dd>The display name of the raider</dd>
    <dt><code>{{ .ToBroadcaster.Login }}</code></dt>
        <dd>The login name of the raided</dd>
    <dt><code>{{ .ToBroadcaster.Name }}</code></dt>
        <dd>The display name of the raided</dd>
    <dt><code>{{ .Viewers }}</code></dt>
        <dd>The number of viewers included in the raid</dd>
</dl>

<h3>Follow Placeholders</h3>
<dl>
    <dt><code>{{ .Follower.Login }}</code></dt>
        <dd>The login name of the follower</dd>
    <dt><code>{{ .Follower.Name }}</code></dt>
        <dd>The display name of the follower</dd>
    <dt><code>{{ .ToBroadcaster.Login }}</code></dt>
        <dd>The login name of the followed</dd>
    <dt><code>{{ .ToBroadcaster.Name }}</code></dt>
        <dd>The display name of the followed</dd>
</dl>

<h3>Cheer Placeholders</h3>
<dl>
    <dt><code>{{ .IsAnonymous }}</code></dt>
        <dd><code>true</code> or <code>false</code>, whether or not the cheer was anonymous</dd>
    <dt><code>{{ .From.Login }}</code></dt>
        <dd>The login name of the cheerer</dd>
    <dt><code>{{ .From.Name }}</code></dt>
        <dd>The display name of the cheerer</dd>
    <dt><code>{{ .Broadcaster.Login }}</code></dt>
        <dd>The login name of the cheered</dd>
    <dt><code>{{ .Broadcaster.Name }}</code></dt>
        <dd>The display name of the cheered</dd>
    <dt><code>{{ .Message }}</code></dt>
        <dd>The message associated with the cheer</dd>
    <dt><code>{{ .Bits }}</code><dt>
        <dd>The number of bits cheered</dd>
</dl>
`;

class Events extends UpdatingControlPanel<banterpb.Config> {
    private _table: HTMLTableElement;

    private _rows: { [key: string]: EventRow } = {};

    constructor(cfg: Cfg) {
        super({ title: 'Events', help, data: cfg });

        this.innerHTML = `
<table></table>
`;

        this._table = this.querySelector('table');
    }

    update(cfg: banterpb.Config) {
        this._table.innerHTML = `
<tr>
    <th>Event</th>
    <th>Enabled</th>
    <th>Text</th>
    <th>Edit</th>
</tr>
`;
        this._rows = {};
        let raid = new EventRow({ name: 'Raid', settings: cfg.channelRaid });
        raid.onsave = (es) => {
            let cfg = this.last.clone();
            cfg.channelRaid = es;
            this.save(cfg);
        }
        this._rows['Raid'] = raid;
        let follow = new EventRow({ name: 'Follow', settings: cfg.channelFollow });
        follow.onsave = (es) => {
            let cfg = this.last.clone();
            cfg.channelFollow = es;
            this.save(cfg);
        }
        this._rows['Follow'] = follow;
        let cheer = new EventRow({ name: 'Cheer', settings: cfg.channelCheer });
        cheer.onsave = (es) => {
            let cfg = this.last.clone();
            cfg.channelCheer = es;
            this.save(cfg);
        }
        this._rows['Cheer'] = cheer;

        Object.keys(this._rows).forEach((key) => {
            let er = this._rows[key];
            er.oncancel = () => this._cancelEditing();
            er.onedit = () => this._setEditing(key);
            this._table.appendChild(er);
        });
    }

    private _cancelEditing() {
        Object.keys(this._rows).forEach((key) => {
            this._rows[key].disabled = false;
        });
    }

    private _setEditing(editKey: string) {
        Object.keys(this._rows).forEach((key) => {
            let er = this._rows[key];
            er.disabled = editKey !== key;
        });
    }
}
customElements.define('banter-events', Events, { extends: 'fieldset' });

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

export { Events };