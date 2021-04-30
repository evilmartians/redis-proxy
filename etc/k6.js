import redis from "k6/x/redis";
import { check } from "k6";

const HOST = __ENV.HOST || "host.docker.internal";
const PORT = __ENV.PORT || 6379;
const PASSWORD = __ENV.PASSWORD || "";
const DBNUM = __ENV.DBNUM || "42";

const client = redis.newClient(`${HOST}:${PORT}`, PASSWORD, Number(DBNUM));

export function setup() {
    redis.do(client, "select", "42");
    redis.do(client, "script", "load", "return redis.call('get', KEYS[1])");
}

export default function () {
    redis.set(client, "hello", "world");
    let res = redis.get(client, "hello");
    check(res, {
        "get hello": (r) => r === "world",
    });

    res = redis.do(
        client,
        "evalsha",
        "4e6d8fc8bb01276962cce5371fa795a7763657ae",
        1,
        "hello"
    );
    check(res, {
        evalsha: (r) => r === "world",
    });

    redis.do(client, "mset", "k1", "v1", "k2", "v2");
    res = redis.do(client, "mget", "k1", "k2");
    check(res, {
        mget: (r) => r[0] === "v1" && r[1] === "v2",
    });
}

export function teardown() {
    redis.del(client, "hello", "k1", "k2");
}
