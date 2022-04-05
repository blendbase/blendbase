const currentWorkflowIssues = [
  {
    id: 1,
    name: "Integrations are costly",
    description: "It can take up to 70% of your development and maintenance efforts"
  },
  {
    id: 2,
    name: "Integrations distract from innovation",
    description:
      "Taking care of app integrations is stealing time and distracting your team from building the differentiating capabilities of your product. "
  },
  {
    id: 3,
    name: "Integrations are never-ending",
    description:
      "With growing numbers of apps and services available on the market, it is becoming more challenging to win and keep users engaged within your product."
  }
];

const newWorkflowBenefits = [
  {
    id: 1,
    name: "Integrate once",
    description: "Build an integration with Blendbase once and benefit from consistent support and updates."
  },
  {
    id: 2,
    name: "Continually expand your coverage",
    description: "Blendbase will be launching new integrations that you can start using with no additional dev cost."
  },
  {
    id: 3,
    name: "Build anything custom on top",
    description: "Build your own integrations on top of Blendbase. It's open source!"
  }
];

export default function Workflows() {
  return (
    <section className="relative">
      {/* <div className="absolute inset-0 bg-gray-100 pointer-events-none mb-16" aria-hidden="true"></div> */}
      {/* <div className="absolute left-0 right-0 m-auto w-px p-px h-20 bg-gray-200 transform -translate-y-1/2"></div> */}

      <div className="relative">
        <div className="relative max-w-6xl mx-auto px-4 sm:px-6">
          <div className="pt-12 md:pt-20">
            {/* Current workflow */}
            <div className="mx-auto text-center pb-4">
              <h3 className="h3 mb-4">Typical SaaS Integrations Dev Workflow</h3>
              <p className="text-xl text-gray-600"></p>
            </div>

            {/* new workflow image */}
            <div className="pb-12 md:pb-16 grid grid-cols-1 lg:grid-cols-3 gap-4">
              <div className="col-span-1 md:col-span-2">
                <img className="w-full mx-auto" src="/images/current-workflow.png" alt="Current workflow" />
              </div>
              <div className="col-span-1">
                <div className="relative">
                  <dl className="space-y-6">
                    {currentWorkflowIssues.map((item) => (
                      <div key={item.id} className="relative">
                        <dt className="mb-1">
                          {/* <div className="absolute flex items-center justify-center h-12 w-12 rounded-md bg-indigo-500 text-white">
                            <item.icon className="h-6 w-6" aria-hidden="true" />
                          </div> */}
                          <p className="text-lg leading-6 font-medium text-gray-900">{item.name}</p>
                        </dt>
                        <dd className="mt-2text-base text-gray-500">{item.description}</dd>
                      </div>
                    ))}
                  </dl>
                </div>
              </div>
            </div>

            <div className="mx-auto text-center pb-4 mt-12">
              <h3 className="h3 mb-4">The Blendbase Unified Approach</h3>
              <p className="text-xl text-gray-600"></p>
            </div>

            {/* new workflow */}
            <div className="pb-12 md:pb-16 grid grid-cols-1 lg:grid-cols-3 gap-4">
              <div className="col-span-1">
                <div className="relative">
                  <dl className="space-y-6">
                    {newWorkflowBenefits.map((item) => (
                      <div key={item.id} className="relative">
                        <dt className="mb-1">
                          {/* <div className="absolute flex items-center justify-center h-12 w-12 rounded-md bg-indigo-500 text-white">
                            <item.icon className="h-6 w-6" aria-hidden="true" />
                          </div> */}
                          <p className="text-lg leading-6 font-medium text-gray-900">{item.name}</p>
                        </dt>
                        <dd className="mt-2text-base text-gray-500">{item.description}</dd>
                      </div>
                    ))}
                  </dl>
                </div>
              </div>
              <div className="col-span-1 md:col-span-2">
                <img className="w-full mx-auto" src="/images/new-workflow.png" alt="Current workflow" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
