import { bus, enumName } from "/bus.js";
import * as buspb from "/pb/bus/bus_pb.js";
import * as banterpb from '/m/banter/pb/banter_pb.js';
import * as twitchpb from '/m/twitch/pb.js';
import { Cfg } from './controller.js';

const TOPIC_REQUEST = enumName(banterpb.BusTopic, banterpb.BusTopic.BANTER_REQUEST);
const TOPIC_TWITCH_REQUEST = enumName(twitchpb.BusTopics, twitchpb.BusTopics.TWITCH_REQUEST);

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
One placeholder is <code>.postCommand</code>, representing all the text that
came after the <code>!command</code>. For example, say you created a command
<code>!said</code> with the text <code>You said "{{ .postCommand }}". That's cool!</code>.
If a user entered <code>!said this is an apple</code>, <code>banter</code>
would send to the channel <code>You said "this is an apple". That's cool!</code>
</p>

<p>
The <code>.sender</code> placeholder has data about the user who ran the command.
You can access specific pieces of data about the user, for example:
<code>{{ .sender.displayName }}</code> will be replaced by sender's display name.
The available details are:
</p>

<dl>
<dt>login</dt>
<dd>The login of the sender</dd>

<dt>displayName</dt>
<dd>The display name of the sender</dd>

<dt>broadcasterType</dt>
<dd><code>affiliate</code> for an affiliate, <code>partner</code> for a partner, and nothing for neither</dd>

<dt>description</dt>
<dd>The user’s description of their channel.</dd>
</dl>

<p>
For example, give a command <code>!bonk</code> with the text
<code>{{ .sender.displayName }} bonks {{ .postCommand }}</code> and the user
<em>AutonomousKoi</em> enters <code>!bonk @SelfDrivingCarp</code>, <code>banter</code>
will send to the channel <code>AutonomousKoi bonks @SelfDrivingCarp</code>.
</p>

<p>
The <code>.original</code> placeholder has details about the original message
the user sent. The available details are:
</p>

<dl>
<dt>text</dt>
<dd>The complete text the user sent to the channel.</dd>

<dt>isMod</dt>
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

class BanterPlayground extends HTMLDialogElement {
    private _cmd: HTMLInputElement;
    private _template: HTMLTextAreaElement;

    private _cfg: Cfg;
    private _sender: twitchpb.User;
    private _banter_id: Number;

    private _render = () => { };
    private _updateSampleCode = () => { };

    constructor(cfg: Cfg) {
        super();

        this._cfg = cfg;

        this.id = 'banter-playground';

        this.innerHTML = `
<h2>Edit Command &#9432;</h2>
<h3>Template</h3>

<div class="grid grid-2-col">

<label for="cmd">Command</label>
<input id="cmd" type="text"/>

<label for="sample">A user enters: <code></code></label>
<input id="sample" type="text"/>

<label for="sender">Sender</label>
<input id="sender" type="text" placeholder="sampleuser" />

</div>

<textarea cols="60" rows="10"></textarea>
<h3>Output</h3>
<div><output></output></div>
<div>
        <button id="save">Save</button>
        <button id="cancel">Cancel</button>
</div>
`;
        this.querySelector('h2').addEventListener('click', () => {
            let nw = window.open('about:blank');
            nw.addEventListener('load', () => {
                let doc = nw.document;
                doc.head.innerHTML = `
    <link href="/main.css" rel="stylesheet">
    <link href="/titillium.css" rel="stylesheet">
    <title>Banter Rendering Help - AutonomousKoi</title>
`;
                doc.body.appendChild(help);
            });
        });
        let sampleCode = this.querySelector('code');
        this._cmd = this.querySelector('input#cmd');
        this._updateSampleCode = () => sampleCode.textContent = this._cmd.value;
        this._cmd.addEventListener('input', () => this._updateSampleCode())

        let sample: HTMLInputElement = this.querySelector('input#sample');

        this._template = this.querySelector('textarea');
        let output = this.querySelector('output');

        this._render = () => {
            let template = this._template.value;
            if (!(template.includes('{{') && template.includes('}}'))) {
                output.value = template;
                return;
            }

            bus.sendAnd(new buspb.BusMessage({
                topic: TOPIC_REQUEST,
                type: banterpb.MessageTypeRequest.BANTER_RENDER_REQ,
                message: new banterpb.BanterRenderRequest({
                    banter: new banterpb.Banter({
                        command: this._cmd.value,
                        text: template,
                    }),
                    original: new twitchpb.event.EventChannelChatMessage({
                        message: new twitchpb.event.ChatMessage({
                            text: this._cmd.value + sample.value ? ' ' + sample.value : '',
                        }),
                        chatter: this._sender,
                    }),
                    sender: this._sender ? this._sender : undefined,
                }).toBinary(),
            })).then((reply) => {
                let resp = banterpb.BanterRenderResponse.fromBinary(reply.message);
                output.value = resp.output;
            }).catch(e => output.value = `ERROR: ${JSON.stringify(e)}`);
        }

        sample.addEventListener('input', () => this._render());
        this._template.addEventListener('input', () => this._render());

        let senderInput: HTMLInputElement = this.querySelector('input#sender');
        senderInput.addEventListener('change', () => {
            if (!senderInput.value) {
                this._sender = undefined;
                return;
            }
            bus.sendAnd(new buspb.BusMessage({
                topic: TOPIC_TWITCH_REQUEST,
                type: twitchpb.requests.MessageTypeRequest.TYPE_REQUEST_GET_USER_REQ,
                message: new twitchpb.requests.GetUserRequest({
                    login: senderInput.value,
                    profile: cfg.last.sendAs,
                }).toBinary(),
            })).then((reply) => {
                let resp = twitchpb.requests.GetUserResponse.fromBinary(reply.message);
                this._sender = resp.user;
                this._render();
            }).catch((e) => {
                output.value = `ERROR: ${JSON.stringify(e)}}`;
            });
        });

        (this.querySelector('button#cancel') as HTMLButtonElement).addEventListener('click', () => this._cancel());
        (this.querySelector('button#save') as HTMLButtonElement).addEventListener('click', () => this._save());
    }

    private _cancel() {
        this.close();
    }

    edit(banter: banterpb.Banter) {
        this._banter_id = banter.id;
        this._cmd.value = banter.command;
        this._template.value = banter.text;
        this.showModal();
        this._updateSampleCode();
        this._render();
    }

    private _save() {
        let newCfg = this._cfg.last.clone();
        let matched = false;
        for (let i = 0; i < newCfg.banters.length; i++) {
            if (newCfg.banters[i].id === this._banter_id) {
                newCfg.banters[i].command = this._cmd.value;
                newCfg.banters[i].text = this._template.value;
                matched = true;
                break;
            }
        }
        if (!matched) {
            newCfg.banters.push(new banterpb.Banter({
                command: this._cmd.value,
                text: this._template.value,
            }));
        }
        this._cfg.save(newCfg);
    }
}
customElements.define('ak-banter-banter-playground', BanterPlayground, { extends: 'dialog' });

export { BanterPlayground, help };