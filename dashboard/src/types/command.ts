type Command = {
    id: number,
    command: string,
    output: string
}

type CommandsResponse = {
    commands: Command[],
}