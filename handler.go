package main

import "sync"

var SETs map[string]string
var SETsMu = sync.RWMutex{}
var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

var handlers = map[string]func([]Value) Value{
	"PING": Ping,
	"SET":  set,
	"GET":  get,
	"HSET": hset,
	"HGET": hget,
}

func Ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ_: STRING, str: "PONG"}
	}
	return Value{typ_: STRING, str: args[0].bulk}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ_: ERROR, str: "ERR wrong number of args for SET"}
	}

	key, val := args[0].bulk, args[1].bulk

	SETsMu.Lock()
	SETs[key] = val
	SETsMu.Unlock()

	return Value{typ_: STRING, str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ_: ERROR, str: "ERR wrong number of args for GET"}
	}

	key := args[0].bulk

	SETsMu.Lock()
	val, ok := SETs[key]
	if !ok {
		return Value{typ_: NULL}
	}

	return Value{typ_: BULK, bulk: val}
}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ_: ERROR, str: "ERR wrong number of args for HSET"}
	}

	h := args[0].bulk
	key := args[1].bulk
	val := args[2].bulk

	HSETsMu.Lock()
	_, ok := HSETs[h]
	if !ok {
		HSETs[h] = map[string]string{}
	}
	HSETs[h][key] = val
	HSETsMu.Unlock()

	return Value{typ_: STRING, str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ_: ERROR, str: "ERR wrong number of args for HSET"}
	}

	h, key := args[0].bulk, args[1].bulk
	HSETsMu.RLock()
	val, ok := HSETs[h][key]
	HSETsMu.RUnlock()

	if !ok {
		return Value{typ_: NULL}
	}

	return Value{typ_: BULK, bulk: val}
}
