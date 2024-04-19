type Agent = {
    id: number,
    name: string,
    username: string,
    displayName: string,
    osUserId: string,
    osUserGroupId: string,
    homeDirectory: string
}

type AgentsResponse = {
    agents: Agent[]
}