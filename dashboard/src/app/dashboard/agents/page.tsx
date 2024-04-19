import { constants } from "@/app/constants"
import { AgentCard } from "@/components/AgentCard";

async function GetAgents() :Promise<Agent[]>{
  const url = `${constants.apiBaseURL}/agents`;
  const response = await fetch(url);
  if(!response.ok) throw new Error("Could not get agents");

  const agents: AgentsResponse = await response.json();
  return agents.agents;
}

export default async function AgentsPage() {
  const agents = await GetAgents();

  return (
    <main className="flex flex-1 flex-col gap-4 p-4 md:gap-8 md:p-8">
      <div className="grid gap-4 md:grid-cols-2 md:gap-8 lg:grid-cols-4">
        {agents.map((agent) => { 
          return <AgentCard key={agent.id} agent={agent}></AgentCard>;
        })}
      </div>
    </main>
  )
}
