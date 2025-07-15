import * as buspb from "/pb/bus/bus_pb.js";
import { bus, enumName } from "/bus.js";
import { requests } from "/m/twitch/pb.js";

let tryImport = (resolve: (t: Twitch) => void) => import('/m/twitch/pb.js')
    .then((mod) => {
        resolve(new Twitch(mod));
    }).catch(() => {
        setTimeout(() => {tryImport(resolve)},500);
    }); 

let twitch = new Promise<Twitch>((resolve) => {
    tryImport(resolve);
});

type modType = typeof import('/m/twitch/pb.js');
type reqType = typeof import('/m/twitch/pb/request_pb.js');

class Twitch {
    private TOPIC_REQUEST = "";
    private static requests: reqType;

    constructor(mod: modType) {
        this.TOPIC_REQUEST = enumName(mod.BusTopics, mod.BusTopics.TWITCH_REQUEST);
        Twitch.requests = mod.requests;
    }

    getUser(id: string): Promise<requests.GetUserResponse> {
        return bus.sendAnd(new buspb.BusMessage({
            topic: this.TOPIC_REQUEST,
            type: requests.MessageTypeRequest.TYPE_REQUEST_GET_USER_REQ,
            message: new requests.GetUserRequest({
                profile: 'selfdrivingcarp',
                login: id,
            }).toBinary(),
        })).then((reply) => {
            return requests.GetUserResponse.fromBinary(reply.message);
        });
    }
}

export { twitch };