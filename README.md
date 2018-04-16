# flyte-hipchat

A HipChat pack for Flyte.

## Bulding

Must have [dep](https://github.com/golang/dep) installed
* Run `dep ensure`
* Run `go test ./...`
* Run `go build`


## Configuration

The plugin is configured using environment variables:

ENV VAR           | Default  |  Description                            | Example                                    
 ---------------- |  ------- |  -------------------------------------- |  ------------------------------------------
FLYTE_API         | -        | The API endpoint to use                 | http://localhost:8080
HIPCHAT_TOKENS    | -        | The API tokens to use, comma separated  | token_abc
DEFAULT_JOIN_ROOM | -        | A room to join by default when launched | 1234
BKP_DIR           | $TMPDIR  | Directory where to backup joined rooms  | /flyte-hipchat

Example `FLYTE_API=http://localhost:8080 HIPCHAT_TOKENS=token_abc DEFAULT_JOIN_ROOM=1234 ./flyte-hipchat`

## Commands

All the events have the same fields as the command input plus error, which is omitted if the command was successful

### SendMessage

    {
        "roomId": "...", // required
        "message": "..." // required
    }

Returned events

`MessageSent`

    {
        "roomId": "...",
        "message": "..."
    }

`SendMessageFailed`

    {
        "roomId": "...",
        "message": "...",
        "error": "..."
    }

### SendNotification

    {
        "roomId": "...",        // required
        "message": "...",       // required
        "messageFormat": "...", // [text|html] default text
        "notify": "...",        // [true|false] default false
        "color": "...",         // [yellow|green|red|purple|gray|random] defaults to yellow
        "from": "..."           // required
    }

Returned events

`NotificationSent`

    {
        "roomId": "...",
        "message": "...",
        "messageFormat": "...",
        "notify": "...",
        "color": "...",
        "from": "..."
    }

`SendNotificationFailed`

    {
        "roomId": "...",
        "message": "...",
        "messageFormat": "...",
        "notify": "...",
        "color": "...",
        "from": "...",
        "error": "..."
    }

### Broadcast

Same as send message, but without room id. Message will be sent to all the rooms that pack has joined.

    {
        "message": "..." // required
    }

Returned events

`BroadcastSent`

    {
        "message": "..."
    }

`BroadcastFailed`

    {
        "message": "...",
        "error": "..."
    }

### JoinRoom

Joins room, pack will start sending `ReceivedMessage` events there's new message in the room

    {
        "roomId": "..." // required
    }

Returned events

`RoomJoined`

    {
        "roomId": "..."
    }

`JoinRoomFailed`

    {
        "roomId": "...",
        "error": "..."
    }

### LeaveRoom

Leaves HipChat room, room will not be monitored for incoming messages.

    {
        "roomId": "..." // required
    }

Returned events

`RoomLeft`

    {
        "roomId": "..."
    }

`LeaveRoomFailed`

    {
        "roomId": "...",
        "error": "..."
    }

## Events 

### ReceivedMessage

	{
        "id": "...",
        "roomId": "...",
        "date": "...",
        "from": {
            "id": "...",
            "name": "...",
            "mentionName": "...",
        },
        "mentions": [
            {
                "id": "...",
                "name": "...",
                "mentionName": "...",
            }
        ],
        "message": "...",
        "messageFormat": "...",
        "type": "..."
    }
