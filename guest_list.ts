import * as banterpb from '/m/banter/pb/banter_pb.js';
import { Cfg } from './controller.js';
import { twitch } from './twitch.js';
import { UpdatingControlPanel } from '/tk.js';

let help = document.createElement('div');
help.innerHTML = `
<p>Guest lists are lists of Twitch users that you can refer to in other parts of Banter,
such as Guest List Commands.</p>`;

class GuestLists extends UpdatingControlPanel<banterpb.Config> {
    private _button: HTMLButtonElement;
    private _table: HTMLTableElement;
    private _edit: GuestListEdit;

    constructor(cfg: Cfg) {
        super({ title: 'Guest Lists', help, data: cfg });

        this.innerHTML = `
<table></table>
<button> + </button>
`;

        this._button = this.querySelector('button');
        this._button.addEventListener('click', () => this._newGuestList());
        this._table = this.querySelector('table');

        this._edit = new GuestListEdit(cfg);
        this.appendChild(this._edit);
    }

    update(cfg: banterpb.Config) {
        let names = Object.keys(cfg.guestLists)
        if (names.length === 0) {
            this._table.textContent = '';
            return;
        }
        this._table.innerHTML = `
<tr>
    <th>Name</th>
    <th>Members</th>
    <th></th>
</tr>
`;
        names = names.toSorted();
        names.forEach((name) => {
            let list = new GuestList(name, cfg.guestLists[name]);
            list.onDelete = () => this._deleteList(name);
            list.onEdit = () => this._edit.editing = name;
            this._table.appendChild(list);
        });
    }

    private _newGuestList() {
        let name = prompt('Guest list name');
        if (this.last.guestLists[name]) {
            alert('name already in use');
            return;
        }
        if (!name) {
            return;
        }
        let cfg = this.last.clone();
        cfg.guestLists[name] = new banterpb.GuestList();
        this.save(cfg);
    }

    private _deleteList(name: string) {
        if (!confirm(`Delete Guest List ${name}?`)) {
            return;
        }
        let cfg = this.last.clone();
        delete cfg.guestLists[name];
        this.save(cfg);
    }
}
customElements.define('banter-guest-lists', GuestLists, { extends: 'fieldset' });

class GuestChip extends HTMLDivElement {
    constructor(guest: banterpb.GuestList_Member, onClick = () => { }) {
        super();

        this.innerHTML = `${guest.login} <button> X </button>`;

        let button = this.querySelector('button');
        button.addEventListener('click', onClick);
    }
}
customElements.define('banter-guest-chip', GuestChip, { extends: 'div' });

class GuestListEdit extends HTMLDialogElement {
    private _guests: HTMLDivElement;
    private _input: HTMLInputElement;

    private _cfg: Cfg;
    private _name = '';
    private _list = new banterpb.GuestList();

    constructor(cfg: Cfg) {
        super();

        this.innerHTML = `
<style>
#guests {
    display: flex;
    flex-wrap: wrap;
}
#guests > div {
    border: solid gray 1px;
    border-radius: 5px;
    padding: 0.5rem;
}
</style>

<div class="flex-column" style="gap: 1rem">
<h3>Manage Guest List Members</h3>
<section>
    <input type="text" size="30" placecholder="@selfdrivingcarp" />
    <button id="add"> + </button>
</section>
<section id="guests"></section>
<section style="align-self: flex-end">
    <button id="save">Save</button>
    <button id="cancel">Cancel</button>
</section>
</div>
`;

        this.close();
        this._input = this.querySelector('input');
        let addGuest = () => this._addGuest(this._input.value);
        this._input.addEventListener('keypress', (e) => {
            let keyCode = e.code || e.key;
            if (keyCode == 'Enter') {
                addGuest();
            }
        });
        let add = this.querySelector('#add');
        add.addEventListener('click', addGuest);

        this._guests = this.querySelector('#guests');

        let save = this.querySelector('#save');
        save.addEventListener('click', () => this._save());
        let cancel = this.querySelector('#cancel');
        cancel.addEventListener('click', () => this.editing = '');

        this.editing = '';

        this._cfg = cfg;
        this._cfg.subscribe((_) => this.editing = '');
    }

    private _save() {
        let cfg = this._cfg.last.clone();
        cfg.guestLists[this._name] = this._list;
        this._cfg.save(cfg);
    }

    set editing(name: string) {
        this._name = name;
        if (!name) {
            this.close();
            this._list = new banterpb.GuestList();
            this._populateChips();
            return;
        }

        if (!(name in this._cfg.last.guestLists)) {
            return;
        }
        this._list = this._cfg.last.guestLists[name].clone();
        this._populateChips();

        this.showModal();
    }

    private _populateChips() {
        this._guests.textContent = '';
        this._list.members.toSorted((a, b) => a.login.localeCompare(b.login))
            .forEach((guest) => {
                let chip = new GuestChip(guest, () => {
                    this._list.members = this._list.members.filter((member) => member.id !== guest.id);
                    this._populateChips();
                });
                this._guests.appendChild(chip);
            });
    }

    private _addGuest(name: string) {
        this._input.value = '';
        if (!name) {
            return;
        }
        twitch.then((twitch) => {
            return twitch.getUser(name);
        }).then((gur) => {
            let guest = new banterpb.GuestList_Member({
                login: gur.login,
                id: gur.user.id,
            })
            if (this._list.members.some((v) => v.id === guest.id)) {
                this._input.value = '';
                return;
            }
            this._list.members.push(guest);
            this._populateChips();
        });
    }
}
customElements.define('ak-guest-list-edit', GuestListEdit, { extends: 'dialog' });

class GuestList extends HTMLTableRowElement {

    onDelete = () => { };
    onEdit = () => { };

    constructor(name: string, l: banterpb.GuestList) {
        super();

        this.innerHTML = `
<td>${name}</td>
<td>${l.members.length}</td>
<td>
    <button id="edit">Edit</button>
    <button id="delete">Delete</button>
</td>
`;

        let del = this.querySelector('#delete');
        del.addEventListener('click', () => this.onDelete());
        let edit = this.querySelector('#edit');
        edit.addEventListener('click', () => this.onEdit());
    }
}
customElements.define('banter-guest-list', GuestList, { extends: 'tr' });

export { GuestLists };