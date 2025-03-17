import * as banterpb from '/m/banter/pb/banter_pb.js';
import { Cfg } from './controller.js';
import { UpdatingControlPanel } from '/tk.js';

let help = document.createElement('div');
help.innerHTML = `
<h3>Custom Commands and Messages</h3>
<p>
You can define custom chat commands and messages to provide information and interactions
to chat. To create a new command, click the <code>+</code> button in the bottom left. A
new empty row will appear. 
</p>

<dl>
    <dt>Command</dt>
    <dd>The command users will type into chat to get this message. It must
        begin with <code>!</code> and must be a single word.
    </dd>

    <dt>Text</dt>
    <dd>The text sent to the chat when the command is activated by a user or
        as a random command. This text can have special placeholder values that
        are replaced before the message is sent to chat. These placeholders are
        described below.
    </dd>

    <dt>Enabled</dt>
    <dd>If not enabled, this command can't be invoked by users, won't be selected
        as a random command, and won't appear in the command list reported by
        <code>!banter</code>.
    </dd>

    <dt>Random</dt>
    <dd>If checked, this command can be randomly selected and sent to the channel.</dd>
</dl>

    <h4>Message Processing</h4>

<p>
When a command is triggered by a user in chat the text will go through processing
before being sent to chat. The text can include special values that will be
replaced with something based on context.
</p>

<p>
The placeholders begin and end with double curly braces: (<code>{{ placholder }}</code>).
One placeholder is <code>.PostCommand</code>, representing all the text that
came after the <code>!command</code>. For example, say you created a command
<code>!said</code> with the text <code>You said "{{ .PostCommand }}". That's cool!</code>.
If a user entered <code>!said this is an apple</code>, <code>banter</code>
would send to the channel <code>You said "this is an apple". That's cool!</code>
</p>

<p>
The <code>.Sender</code> placeholder has data about the user who ran the command.
You can access specific pieces of data about the user, for example:
<code>{{ .Sender.DisplayName }}</code> will be replaced by sender's display name.
The available details are:
</p>

<dl>
<dt>Login</dt>
<dd>The login of the sender</dd>

<dt>DisplayName</dt>
<dd>The display name of the sender</dd>

<dt>BroadcasterType</dt>
<dd><code>affiliate</code> for an affiliate, <code>partner</code> for a partner, and nothing for neither</dd>

<dt>Description</dt>
<dd>The user’s description of their channel.</dd>
</dl>

<p>
For example, give a command <code>!bonk</code> with the text
<code>{{ .Sender.DisplayName }} bonks {{ .PostCommand }}</code> and the user
<em>AutonomousKoi</em> enters <code>!bonk @SelfDrivingCarp</code>, <code>banter</code>
will send to the channel <code>AutonomousKoi bonks @SelfDrivingCarp</code>.
</p>

<p>
The <code>.Original</code> placeholder has details about the original message
the user sent. It has details like <code>.Sender</code>. The available details are:
</p>

<dl>
<dt>Text</dt>
<dd>The complete text the user sent to the channel.</dd>

<dt>IsMod</dt>
<dd>Whether or not the sender is a channel mod.</dd>
</dl>

<p>
It's possible to create invalid placeholders. If you save an invalid message
<code>banter</code> will save it but it won't appear in chat. When an invalid
command is run an error with details will be written to the logs.
</p>

<p>
Under the hood <code>banter</code> is using Go's
<a href="https://pkg.go.dev/text/template">text/template</a> for text processing.
All of that package's functionality is available if you're adventerous. There's
no limit to the length of the text <code>banter</code> will save for a command
but there is a limit to the length of a valid chat message. The full details
of the data available to the template are in the source for
<a href="https://github.com/autonomouskoi/twitch/blob/main/twitch.pb.go">Sender</a>
(a <code>User</code> struct) and
<a href="https://github.com/autonomouskoi/twitch/blob/main/chat.pb.go">Original</a>
(a <code>TwitchChatEventMessageIn</code> struct).
</p>
`;

class BanterMessages extends UpdatingControlPanel<banterpb.Config> {
    private _table: HTMLTableElement;
    private _button: HTMLButtonElement;

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
            row.oncancel = () => this._cancelEditing();
            row.onedit = () => this._setEditing(i);
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
    }

    private _handleNew() {
        let row = new BanterRow();
        row.banter = new banterpb.Banter();
        row.oncancel = () => {
            this._cancelEditing();
        };
        row.onsave = (banter) => {
            let cfg = this.last.clone();
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