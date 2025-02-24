import { SectionHelp } from "/help.js";
import * as banterpb from '/m/banter/pb/banter_pb.js';
import { Cfg } from './controller.js';


class GuestListCommands extends HTMLFieldSetElement {
    private _button: HTMLButtonElement;
    private _table: HTMLTableElement;
    private _edit: GuestListCommandEdit;

    private _cfg: Cfg;

    constructor(cfg: Cfg) {
        super();

        this.innerHTML = `
<legend>Guest List Commands &#9432;</legend>

<div id="help"></div>

<table></table>
<button> + </button>
`;

        this._setHelp();
        this._cfg = cfg;
        this._button = this.querySelector('button');
        this._button.addEventListener('click', () => this._newCmd());
        this._table = this.querySelector('table');

        this._edit = new GuestListCommandEdit(this._cfg);
        this.appendChild(this._edit);

        this._cfg.subscribe((newCfg) => this.config = newCfg);
        this.config = this._cfg.last;
    }

    set config(cfg: banterpb.Config) {
        let names = Object.keys(cfg.guestListCommands);
        if (names.length === 0) {
            this._table.textContent = '';
            return;
        }
        this._table.innerHTML = `
    <tr>
        <th>Name</th>
        <th>Lists</th>
        <th></th>
    </tr>
    `;
        names = names.toSorted();
        names.forEach((name) => {
            let cmd = new GuestListCommand(name, cfg.guestListCommands[name]);
            cmd.onDelete = () => this._deleteCmd(name);
            cmd.onEdit = () => this._edit.editing = name;
            this._table.appendChild(cmd);
        });
    }

    private _newCmd() {
        let name = prompt('Command name (not the command itself)');
        if (this._cfg.last.guestListCommands[name]) {
            alert('name already in use');
            return;
        }
        if (!name) {
            return;
        }
        let cfg = this._cfg.last.clone();
        cfg.guestListCommands[name] = new banterpb.GuestListCommand();
        this._cfg.save(cfg);
    }

    private _deleteCmd(name: string) {
        if (!confirm(`Delete Guest List Command ${name}?`)) {
            return;
        }
        let cfg = this._cfg.last.clone();
        delete cfg.guestListCommands[name];
        this._cfg.save(cfg);
    }

    private _setHelp() {
        let helpToggle = this.querySelector('legend');

        let helpHTML = document.createElement('div');
        helpHTML.innerHTML = `
<p>
The first time someone chats during a session Banter will check to see if they are on any Guest Lists. Each
Guest List they are on can have associated commands. Each of those commands is triggered for the user.
</p>

<p>
The command text goes through template processing before being sent to Twitch. For example, if
<code>selfdrivingcarp</code> is on a guest list for the command:

<blockquote>
    <code>Hey {{ .DisplayName }}, whatup!?!?!</code>
</blockquote>

then

<blockquote>
    <code>Hey SelfDrivingCarp, whatup!?!?!</code>
</blockquote>

will be sent to chat.
</p>

<p>
Things in double curly braces (<code>{{ }}</code>) are interpeted by Banter and usually have user properties.
You can find the available properties
<a href="https://github.com/autonomouskoi/twitch/blob/main/twitch.pb.go#L89" target="_blank">with the source</a>. Note that not
all properties will have values for all users.
</p>
`;
        let help = SectionHelp(helpToggle, helpHTML);
        let helpPlaceholder = this.querySelector('#help');
        this.replaceChild(help, helpPlaceholder);
    }
}
customElements.define('banter-guest-list-commands', GuestListCommands, { extends: 'fieldset' });

class GuestListCommand extends HTMLTableRowElement {

    onDelete = () => { };
    onEdit = () => { };

    constructor(name: string, c: banterpb.GuestListCommand) {
        super();

        this.innerHTML = `
<td>${name}</td>
<td>${c.guestListNames.toSorted().join(", ")}</td>
<td>
    <button id="edit">Edit</button>
    <button id="delete">Delete</button>
</td>
`;

        this.querySelector('#delete').addEventListener('click', () => this.onDelete());
        this.querySelector('#edit').addEventListener('click', () => this.onEdit());
    }
}
customElements.define('banter-guest-list-command', GuestListCommand, { extends: 'tr' });

class GuestListCommandEdit extends HTMLDialogElement {
    private _lists: HTMLDivElement;
    private _input: HTMLInputElement;

    private _cfg: Cfg;
    private _name = '';

    constructor(cfg: Cfg) {
        super();

        this.innerHTML = `
<div class="flex-column" style="gap: 1rem">
<section>
    <h3>Chat message:</h3>
    <input type="text" size="50" placeholder="Hi, {{ .DisplayName }}!" />
</section>
<section id="lists" class="grid grid-2-col"></section>
<section style="align-self: flex-end">
    <button id="save">Save</button>
    <button id="cancel">Cancel</button>
</section>
</div>
`;
        this.close();
        this._input = this.querySelector('input');
        this._lists = this.querySelector('#lists');

        this.querySelector('#cancel').addEventListener('click', () => this.editing = '');
        this.querySelector('#save').addEventListener('click', () => this._save());

        this._cfg = cfg;
        this._cfg.subscribe((_) => this.editing = '');
    }

    private _save() {
        let cfg = this._cfg.last.clone();
        let cmd = new banterpb.GuestListCommand({
            command: this._input.value,
            guestListNames: Object.keys(cfg.guestLists).filter((name) => {
                let elem: HTMLInputElement = this.querySelector('#' + name);
                if (!elem) { return false }
                return elem.checked;
            }),
        });
        cfg.guestListCommands[this._name] = cmd;
        this._cfg.save(cfg);
    }

    set editing(name: string) {
        this._name = name;
        if (!name) {
            this.close();
            this._display();
            return;
        }
        this._display();
        this.showModal();
    }

    private _display() {
        this._lists.innerHTML = '<h3 style="grid-column: 1 / span 2">Trigger for members of:</h3>';
        if (!this._name) {
            return;
        }
        this._input.value = this._cfg.last.guestListCommands[this._name].command;
        let cmdLists = this._cfg.last.guestListCommands[this._name].guestListNames;
        Object.keys(this._cfg.last.guestLists).toSorted()
            .forEach((name) => {
                let label = document.createElement('label');
                label.innerText = name;
                label.htmlFor = name;
                this._lists.appendChild(label);

                let check = document.createElement('input');
                check.type = 'checkbox';
                check.id = name;
                check.checked = cmdLists.includes(name);
                this._lists.appendChild(check);
            });
    }
}
customElements.define('banter-guest-list-command-edit', GuestListCommandEdit, { extends: 'dialog' });

export { GuestListCommands };