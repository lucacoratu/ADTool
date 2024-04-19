import { FC } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "./ui/card";
import Link from "next/link";

type AgentCardProps = {
    agent: Agent
}

export const AgentCard: FC<AgentCardProps> = ({agent}): JSX.Element => {
    return (
        <Link href={`/dashboard/agents/${agent.id}`}>
            <Card className="">
                <CardHeader>
                    <CardTitle>{agent.name ? agent.name : "No name"}</CardTitle>
                    <CardDescription>Id: {agent.id}</CardDescription>
                </CardHeader>
                <CardContent className="flex flex-col gap-2 text-sm">
                    <div className="flex flex-col">
                        <div className="flex flex-row justify-between">
                            <p>OS User:</p> 
                            <p>{agent.username}</p>
                        </div>
                        <p className="text-sm text-muted-foreground">{agent.osUserId}</p>
                    </div>
                    <div className="flex flex-row justify-between">
                        <p>OS Group:</p>
                        <p>{agent.osUserGroupId}</p>
                    </div>
                    <div className="flex flex-row justify-between">
                        <p>Home Directory:</p>
                        <p>{agent.homeDirectory}</p>
                    </div>
                </CardContent>
            </Card>
        </Link>
    );
}