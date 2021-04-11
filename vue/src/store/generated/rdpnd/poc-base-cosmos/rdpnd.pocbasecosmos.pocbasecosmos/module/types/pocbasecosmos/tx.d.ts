import { Reader, Writer } from "protobufjs/minimal";
import { Coin } from "../cosmos/base/v1beta1/coin";
export declare const protobufPackage = "rdpnd.pocbasecosmos.pocbasecosmos";
export interface MsgAdminSpendCommunityPool {
    initiator: string;
    toAddress: string;
    coins: Coin[];
}
export interface MsgAdminSpendResponse {
}
export declare const MsgAdminSpendCommunityPool: {
    encode(message: MsgAdminSpendCommunityPool, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): MsgAdminSpendCommunityPool;
    fromJSON(object: any): MsgAdminSpendCommunityPool;
    toJSON(message: MsgAdminSpendCommunityPool): unknown;
    fromPartial(object: DeepPartial<MsgAdminSpendCommunityPool>): MsgAdminSpendCommunityPool;
};
export declare const MsgAdminSpendResponse: {
    encode(_: MsgAdminSpendResponse, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): MsgAdminSpendResponse;
    fromJSON(_: any): MsgAdminSpendResponse;
    toJSON(_: MsgAdminSpendResponse): unknown;
    fromPartial(_: DeepPartial<MsgAdminSpendResponse>): MsgAdminSpendResponse;
};
/** Msg defines the Msg service. */
export interface Msg {
    /** this line is used by starport scaffolding # proto/tx/rpc */
    AdminSpendCommunityPool(request: MsgAdminSpendCommunityPool): Promise<MsgAdminSpendResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    AdminSpendCommunityPool(request: MsgAdminSpendCommunityPool): Promise<MsgAdminSpendResponse>;
}
interface Rpc {
    request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}
declare type Builtin = Date | Function | Uint8Array | string | number | undefined;
export declare type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
