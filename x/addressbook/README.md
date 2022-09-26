# Addressbook Module Specification

## Overview

A module to store addresses from different networks related to given cudos address. It is useful for example when given cudos address should 
receive rewards on different network, we can use this module to set and then lookup the other network address.
When creating address, it is associated with network and label, so you can store multiple address for same network related to the same cudos address.

## Module Interface

#### Transaction

| Command                               | Description                                                                         |
| --------------------------------------|-------------------------------------------------------------------------------------|
| [`create-address`](#create-address)   | Creates address entry associated with the sender cudos address, network and label.  |
| [`delete-address`](#delete-address)   | Deletes entry in the addressbook.                                                   |
| [`update-address`](#update-address)   | Updates the address of existing entry in the addressbook.                           |

#### Query

| Command                               | Description                                |
| --------------------------------------|--------------------------------------------|
| [`list-addresses`](#list-addresses)   | Lists all entries from the addressbook.    |
| [`show-address`](#show-address)       | Show address by creator, network and label |


## Object types:

### Address

```go
type Address struct {
	Network string `protobuf:"bytes,1,opt,name=network,proto3" json:"network,omitempty"`
	Label   string `protobuf:"bytes,2,opt,name=label,proto3" json:"label,omitempty"`
	Value   string `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
	Creator string `protobuf:"bytes,4,opt,name=creator,proto3" json:"creator,omitempty"`
}
```

## Full commands info

### Transactions

### `create-address`

> Creates entry in the addressbook associated with creator address, network and label.

- arguments:
  - `network` `string` `Network where the address belongs` `required: true`
  - `label` `string` `Label for the address` `required: true`
  - `value` `string` `Address` `required: true`
- flags:
  none

```bash
cudos-noded tx addressbook create-address BTC 1@TestDenom 1PUo5nwAqp4f2vyzR4ZmRfE3bwyHo8APa5 --keyring-backend=<keyring> --chain-id=<chain-id> --gas=auto --gas-adjustment=1.3 --gas-prices=5000000000000acudos --from=<from-key>
```

### `delete-address`

> Delete addressbook entry.

- arguments:
  - `network` `string` `Network where the address belongs` `required: true`
  - `label` `string` `Label for the address` `required: true`
- flags:
  none

```bash
cudos-noded tx addressbook delete-address BTC 1@TestDenom --keyring-backend=<keyring> --chain-id=<chain-id> --gas=auto --gas-adjustment=1.3 --gas-prices=5000000000000acudos --from=<from-key>
```

### `update-address`

> Update addressbook entry.

- arguments:
  - `network` `string` `Network where the address belongs` `required: true`
  - `label` `string` `Label for the address` `required: true`
  - `value` `string` `Address` `required: true`
- flags:
  none

```bash
cudos-noded tx addressbook update-address BTC 1@TestDenom 1LhrDpLmq6yTgZQ6xPMUsMhXuEHdG3evSq --keyring-backend=<keyring> --chain-id=<chain-id> --gas=auto --gas-adjustment=1.3 --gas-prices=5000000000000acudos --from=<from-key>
```

### Queries

### `list-addresses`

> List all addressbook entries.

- arguments:
  - none
- flags:
  - none

```bash
cudos-noded query addressbook list-addresses
```

### `show-address`

> Show address by creator, network and label.

- arguments:
  - none
- flags:
  - none

```bash
cudos-noded query addressbook show-address cudos1a326k254fukx9jlp0h3fwcr2ymjgludzum67dv BTC 1@TestDenom
```
