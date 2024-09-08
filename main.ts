import { bus, enumName } from "/bus.js";
import * as buspb from "/pb/bus/bus_pb.js";

import * as banterpb from "/m/banter/pb/banter_pb.js";
import * as config from "./config.js";

const TOPIC_BANTER_REQUEST = enumName(banterpb.BusTopic, banterpb.BusTopic.BANTER_REQUEST);

function start(mainContainer: HTMLElement) {
    document.querySelector("title").innerText = 'Banter';

    let cfgElem = new config.Config();

    bus.waitForTopic(TOPIC_BANTER_REQUEST, 5000)
        .then(() => {
            let msg = new buspb.BusMessage();
            msg.topic = TOPIC_BANTER_REQUEST;
            msg.type = banterpb.MessageTypeRequest.CONFIG_GET_REQ;
            msg.message = new banterpb.ConfigGetRequest().toBinary();
            bus.sendWithReply(msg, (reply: buspb.BusMessage) => {
                if (reply.error) {
                    throw reply.error;
                }
                let cgr = banterpb.ConfigGetResponse.fromBinary(reply.message);
                cfgElem.config = cgr.config;
                mainContainer.appendChild(cfgElem);
            });
        });

}

export { start };