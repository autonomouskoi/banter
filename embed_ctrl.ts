import { bus, enumName } from "/bus.js";

import * as banterpb from "/m/banter/pb/banter_pb.js";
import { Cfg } from './controller.js';
import { Random } from "./random.js";
import { BanterMessages } from "./messages.js";
import { Events } from "./events.js";
import { GuestLists } from "./guest_list.js";
import { GuestListCommands } from "./guest_list_commands.js";

const TOPIC_REQUEST = enumName(banterpb.BusTopic, banterpb.BusTopic.BANTER_REQUEST);

function start(mainContainer: HTMLElement) {
    let cfg = new Cfg();

    mainContainer.classList.add('flex-column');
    mainContainer.style.setProperty('gap', '1rem');

    bus.waitForTopic(TOPIC_REQUEST, 5000)
        .then(() => {
            mainContainer.appendChild(new BanterMessages(cfg));
            mainContainer.appendChild(new Random(cfg));
            mainContainer.appendChild(new Events(cfg));
            mainContainer.appendChild(new GuestLists(cfg));
            mainContainer.appendChild(new GuestListCommands(cfg));
            cfg.refresh();
        });
}

export { start };