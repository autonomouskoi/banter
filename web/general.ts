import * as banterpb from '/m/banter/pb/banter_pb.js';
import { Cfg } from './controller.js';
import { UpdatingControlPanel } from '/tk.js';
import { ProfileSelector } from '/m/twitch/profiles.js';

let help = document.createElement('div');
help.innerHTML = `
<dl>
    <dt>Twitch Profile</dt>
    <dd>Select the Twitch profile to use when sending messages.</dd>
</dl>
`;

class General extends UpdatingControlPanel<banterpb.Config> {
    private _sendAs: ProfileSelector;
    private _sendTo: ProfileSelector;
    constructor(cfg: Cfg) {
        super({ title: 'General', help, data: cfg });

        let sendAsTitle = 'Which Twitch profile to send chat messages as';
        let sendToTitle = 'Which channel to send messages to';
        this.innerHTML = `
<div class="grid grid-2-col">
    <label for="send-as" title="${sendAsTitle}">Send As</label>
    <label for="send-to" title="${sendToTitle}">To Channel</label>
</div>
`;

        this._sendAs = new ProfileSelector();
        this._sendAs.id = 'send-as';
        this._sendAs.title = sendAsTitle;
        this._sendAs.addEventListener('change', () => {
            let newCfg = cfg.last.clone();
            newCfg.sendAs = this._sendAs.value;
            this.save(newCfg);
        });

        this._sendTo = new ProfileSelector();
        this._sendTo.id = 'send-to';
        this._sendTo.title = sendToTitle;
        this._sendTo.addEventListener('change', () => {
            let newCfg = cfg.last.clone();
            newCfg.sendTo = this._sendTo.value;
            this.save(newCfg);
        });

        this.querySelector('label[for="send-as"]').after(this._sendAs);
        this.querySelector('label[for="send-to"]').after(this._sendTo);
    }

    update(cfg: banterpb.Config) {
        this._sendAs.selected = cfg.sendAs;
        this._sendTo.selected = cfg.sendTo;
    }
}
customElements.define('ak-banter-general', General, { extends: 'fieldset' });

export { General };