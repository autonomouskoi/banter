import * as banterpb from '/m/banter/pb/banter_pb.js';
import { Cfg } from './controller.js';

class Random extends HTMLFieldSetElement {
    private _cfg: Cfg;
    private _interval: HTMLInputElement;
    private _cooldown: HTMLInputElement;

    constructor(cfg: Cfg) {
        super();
        this._cfg = cfg;

        this.innerHTML = `
<legend>Random Command Configuration</legend>

<label for="input-interval-seconds">Command Interval (seconds)</label>
<input id="input-interval-seconds" type="text"
    inputmode="numeric" pattern="\\d+" size="4" placeholder="300" value="300"
    title="How long between random chat commands""
/>

<label for="input-cooldown-seconds">Command Cooldown (seconds)</label>
<input id="input-cooldown-seconds" type="type"
    inputmode="numeric" pattern="\\d+" size="4" placeholder="900" value="900"
    title="Minimum time before repeating a given random command"
/>
`;
        this.classList.add('grid', 'grid-2-col');
        this._interval = this.querySelector('#input-interval-seconds');
        this._interval.addEventListener('change', () => this._save());
        this._cooldown = this.querySelector('#input-cooldown-seconds');
        this._cooldown.addEventListener('change', () => this._save());

        cfg.subscribe((newCfg) => this.config = newCfg);
        this.config = cfg.last;
    }

    set config(cfg: banterpb.Config) {
        this._interval.value = cfg.intervalSeconds.toString();
        this._cooldown.value = cfg.cooldownSeconds.toString();
    }

    private _save() {
        let interval = parseInt(this._interval.value);
        if (interval < 1) {
            return;
        }
        let cooldown = parseInt(this._cooldown.value);
        if (cooldown < 0) {
            return;
        }
        let cfg = this._cfg.last.clone();
        cfg.intervalSeconds = interval;
        cfg.cooldownSeconds = cooldown;
        this._cfg.save(cfg);
    }
}
customElements.define('banter-random', Random, {extends: 'fieldset'});

export { Random };