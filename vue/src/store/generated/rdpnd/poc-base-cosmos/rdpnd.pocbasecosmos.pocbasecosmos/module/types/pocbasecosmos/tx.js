/* eslint-disable */
import { Reader, Writer } from "protobufjs/minimal";
import { Coin } from "../cosmos/base/v1beta1/coin";
export const protobufPackage = "rdpnd.pocbasecosmos.pocbasecosmos";
const baseMsgAdminSpendCommunityPool = { initiator: "", toAddress: "" };
export const MsgAdminSpendCommunityPool = {
    encode(message, writer = Writer.create()) {
        if (message.initiator !== "") {
            writer.uint32(10).string(message.initiator);
        }
        if (message.toAddress !== "") {
            writer.uint32(18).string(message.toAddress);
        }
        for (const v of message.coins) {
            Coin.encode(v, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseMsgAdminSpendCommunityPool,
        };
        message.coins = [];
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.initiator = reader.string();
                    break;
                case 2:
                    message.toAddress = reader.string();
                    break;
                case 3:
                    message.coins.push(Coin.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = {
            ...baseMsgAdminSpendCommunityPool,
        };
        message.coins = [];
        if (object.initiator !== undefined && object.initiator !== null) {
            message.initiator = String(object.initiator);
        }
        else {
            message.initiator = "";
        }
        if (object.toAddress !== undefined && object.toAddress !== null) {
            message.toAddress = String(object.toAddress);
        }
        else {
            message.toAddress = "";
        }
        if (object.coins !== undefined && object.coins !== null) {
            for (const e of object.coins) {
                message.coins.push(Coin.fromJSON(e));
            }
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.initiator !== undefined && (obj.initiator = message.initiator);
        message.toAddress !== undefined && (obj.toAddress = message.toAddress);
        if (message.coins) {
            obj.coins = message.coins.map((e) => (e ? Coin.toJSON(e) : undefined));
        }
        else {
            obj.coins = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseMsgAdminSpendCommunityPool,
        };
        message.coins = [];
        if (object.initiator !== undefined && object.initiator !== null) {
            message.initiator = object.initiator;
        }
        else {
            message.initiator = "";
        }
        if (object.toAddress !== undefined && object.toAddress !== null) {
            message.toAddress = object.toAddress;
        }
        else {
            message.toAddress = "";
        }
        if (object.coins !== undefined && object.coins !== null) {
            for (const e of object.coins) {
                message.coins.push(Coin.fromPartial(e));
            }
        }
        return message;
    },
};
const baseMsgAdminSpendResponse = {};
export const MsgAdminSpendResponse = {
    encode(_, writer = Writer.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseMsgAdminSpendResponse };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(_) {
        const message = { ...baseMsgAdminSpendResponse };
        return message;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = { ...baseMsgAdminSpendResponse };
        return message;
    },
};
export class MsgClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
    }
    AdminSpendCommunityPool(request) {
        const data = MsgAdminSpendCommunityPool.encode(request).finish();
        const promise = this.rpc.request("rdpnd.pocbasecosmos.pocbasecosmos.Msg", "AdminSpendCommunityPool", data);
        return promise.then((data) => MsgAdminSpendResponse.decode(new Reader(data)));
    }
}
