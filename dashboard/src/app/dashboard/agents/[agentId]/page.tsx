"use client";

import { constants } from "@/app/constants"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import Link from "next/link"
import { ArrowUpRight } from "lucide-react"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Dialog, DialogContent, DialogDescription, DialogTitle, DialogTrigger } from "@/components/ui/dialog"
import { FormEvent } from "react"


async function GetAgentCommands(agentId: string) {
    const URL = `${constants.apiBaseURL}/agents/${agentId}/cmd`;
    const resp = await fetch(URL);
    if(!resp.ok) throw new Error("Could not get commands");
    const commands: CommandsResponse = await resp.json();
    return commands.commands;
}

export default async function AgentPage({ params }: { params: { agentId: string } }) {
    const agentId: string = params.agentId;
    const commands = await GetAgentCommands(agentId);

    async function onSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()
        
        const URL = `${constants.apiBaseURL}/agents/${agentId}/cmd`;

        const formData = new FormData(event.currentTarget)
        const bodyData = {"command": formData.get("cmd")};

        const response = await fetch(URL, {
          method: 'POST',
          body: JSON.stringify(bodyData),
        })
    }

    async function onSubmitRecurring(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()
        
        const URL = `${constants.apiBaseURL}/agents/${agentId}/reccmd`;

        const formData = new FormData(event.currentTarget)
        const bodyData = {"command": formData.get("cmd"), "interval": Number(formData.get("interval"))};

        const response = await fetch(URL, {
          method: 'POST',
          body: JSON.stringify(bodyData),
        });
    }

  return (
    <main className="flex flex-1 flex-col gap-4 p-4 md:gap-8 md:p-8">
        <div className="grid gap-4 md:grid-cols-2 md:gap-8 lg:grid-cols-3">
            <Card className="xl:col-span-2">
                <CardHeader className="flex flex-row items-center">
                <div className="grid gap-2">
                    <CardTitle>Commands</CardTitle>
                    <CardDescription>
                    Commands run by agent.
                    </CardDescription>
                </div>
                <Button asChild size="sm" className="ml-auto gap-1">
                    <Link href="#">
                    View All
                    <ArrowUpRight className="h-4 w-4" />
                    </Link>
                </Button>
                </CardHeader>
                <CardContent>
                    <Table>
                        <TableHeader>
                        <TableRow>
                            <TableHead className="text-left">ID</TableHead>
                            <TableHead className="text-left">Command</TableHead>
                            <TableHead className="text-right">Output</TableHead>
                        </TableRow>
                        </TableHeader>
                        <TableBody>
                            {commands.map((command) => {
                                return (
                                    <TableRow key={command.id}>
                                        <TableCell className="text-left">
                                            {command.id}
                                        </TableCell>
                                        <TableCell className="md:table-cell text-left">
                                            {command.command}
                                        </TableCell>
                                        <TableCell className="sm:table-cell text-right">
                                            {
                                                command.output === "" ? 
                                                <Badge className="text-xs" variant="destructive" >Unavailable</Badge> 
                                                : 
                                                <Dialog>
                                                    <DialogTrigger>
                                                        <Badge className="text-xs" variant="default">Available</Badge>
                                                    </DialogTrigger>
                                                    <DialogContent className="min-w-fit">
                                                        <DialogTitle>Command Output</DialogTitle>
                                                        <DialogDescription className="min-w-fit">{command.output}</DialogDescription>
                                                    </DialogContent>
                                                </Dialog>
                                            }
                                        </TableCell>
                                    </TableRow>     
                                );
                            })}
                        </TableBody>
                    </Table>
                </CardContent>
            </Card>
            <div className="flex flex-col gap-2">
                <Card className="max-h-fit h-fit">
                    <CardHeader>
                        <CardTitle>New Command</CardTitle>
                        <CardDescription>Execute a new command on the agent</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <form onSubmit={onSubmit}>
                            <Label htmlFor="cmd">Command</Label>
                            <Input id="cmd" name="cmd" placeholder="Insert your command..."></Input>
                            <Button className="mt-3 w-20 h-8">Send</Button>
                        </form>
                    </CardContent>
                </Card>
                <Card className="max-h-fit h-fit">
                    <CardHeader>
                        <CardTitle>New Recurring Command</CardTitle>
                        <CardDescription>Execute a new recurring command on the agent</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <form onSubmit={onSubmitRecurring}>
                            <Label htmlFor="cmd">Command</Label>
                            <Input id="cmd" name="cmd" placeholder="Insert your command..."></Input>
                            <Label htmlFor="interval">Interval</Label>
                            <Input type="number" id="interval" name="interval" placeholder="Insert interval..."></Input>
                            <Button className="mt-3 w-20 h-8">Send</Button>
                        </form>
                    </CardContent>
                </Card>
            </div>
        </div>
    </main>
  )
}