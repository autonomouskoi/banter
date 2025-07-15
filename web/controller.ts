import { bus, enumName } from "/bus.js";
import * as buspb from "/pb/bus/bus_pb.js";
import * as banterpb from '/m/banter/pb/banter_pb.js';
import { ValueUpdater } from "/vu.js";

const TOPIC_REQUEST = enumName(banterpb.BusTopic, banterpb.BusTopic.BANTER_REQUEST);
const TOPIC_COMMAND = enumName(banterpb.BusTopic, banterpb.BusTopic.BANTER_COMMAND);

class Cfg extends ValueUpdater<banterpb.Config> {
    constructor() {
        super(new banterpb.Config());
    }

    refresh() {
        bus.sendAnd(new buspb.BusMessage({
            topic: TOPIC_REQUEST,
            type: banterpb.MessageTypeRequest.CONFIG_GET_REQ,
            message: new banterpb.ConfigGetRequest().toBinary(),
        })).then((reply) => {
            let cgResp = banterpb.ConfigGetResponse.fromBinary(reply.message);
            this.update(cgResp.config);
        });
    }

    save(cfg: banterpb.Config): Promise<void> {
        let csr = new banterpb.ConfigSetRequest();
        csr.config = cfg;
        return bus.sendAnd(new buspb.BusMessage({
            topic: TOPIC_COMMAND,
            type: banterpb.MessageTypeCommand.CONFIG_SET_REQ,
            message: csr.toBinary(),
        })).then((reply) => {
            let csResp = banterpb.ConfigSetResponse.fromBinary(reply.message);
            this.update(csResp.config);
        });
    }
}
export { Cfg };