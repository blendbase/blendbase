import React, { useState, useRef, useEffect } from "react";
import Transition from "../utils/Transition";

const IconArrowRight = () => (
  <svg
    className="w-3 h-3 fill-current"
    width="16"
    height="16"
    viewBox="0 0 16 16"
    xmlns="http://www.w3.org/2000/svg"
    fill="current"
  >
    <line x1="12.8008" y1="7.48596" x2="2.17064" y2="7.48596" stroke="black" strokeWidth="2" />
    <line x1="12.789" y1="7.54658" x2="7.13213" y2="1.88973" stroke="black" strokeWidth="2" />
    <line x1="13.4962" y1="6.83955" x2="7.83931" y2="12.4964" stroke="black" strokeWidth="2" />
  </svg>
);

const features = [
  {
    id: 1,
    name: "Single CRM API",
    description: "GraphQL API to query any customer CRM system, including Salesforce, Hubspot, and more coming soon.",
    imageUrl: "/images/features-crm-api.png"
  },
  {
    id: 2,
    name: "Connect API",
    description: "API to build an Integrations page in your SaaS app that handles credential setting and auth flows.",
    imageUrl: "/images/features-connect-api.png"
  },
  {
    id: 3,
    name: "Sample React App",
    description: "Sample React app showcasing the integration with the Connect API and the CRM API",
    imageUrl: "/images/features-sample-app.png"
  }
];

function FeaturesHome() {
  const [tab, setTab] = useState(1);

  const tabs = useRef(null);

  const heightFix = () => {
    if (tabs.current.children[tab - 1]) {
      tabs.current.style.height = tabs.current.children[tab - 1].offsetHeight + "px";
    }
  };

  useEffect(() => {
    heightFix();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [tab]);

  return (
    <section className="relative" id="features">
      {/* Section background (needs .relative class on parent and next sibling elements) */}
      <div className="absolute inset-0 bg-gray-100 pointer-events-none -mb-16" aria-hidden="true"></div>
      <div className="absolute left-0 right-0 m-auto w-px p-px h-278 bg-gray-200 transform -translate-y-1/2"></div>

      <div className="relative max-w-6xl mx-auto px-4 sm:px-6">
        <div className="pt-12 md:pt-20">
          {/* Section content */}
          <div className="md:grid md:grid-cols-12 md:gap-6">
            {/* Content */}
            <div className="max-w-xl md:max-w-none md:w-full mx-auto md:col-span-7 lg:col-span-6">
              {/* Tabs buttons */}
              <div className="mb-8 md:mb-0">
                {features.map((feature) => (
                  <button
                    className={`flex items-center text-lg p-5 rounded border transition duration-300 ease-in-out mb-3 ${
                      tab !== feature.id
                        ? "bg-white shadow-md border-gray-200 hover:shadow-lg"
                        : "bg-brand-50 border-transparent"
                    }`}
                    onClick={(e) => {
                      e.preventDefault();
                      setTab(feature.id);
                    }}
                    key={feature.id}
                  >
                    <div>
                      <div className="font-bold font-grotesk leading-snug tracking-tight mb-1 text-brand-600">
                        {feature.name}
                      </div>
                      <div className="text-slate-600">{feature.description}</div>
                    </div>
                    <div className="flex text-slate-300 justify-center items-center w-8 h-8 bg-white rounded-full shadow shrink-0 ml-3">
                      <IconArrowRight />
                    </div>
                  </button>
                ))}
              </div>
            </div>

            {/* Tabs items */}
            <div
              className="max-w-xl md:max-w-none md:w-full mx-auto md:col-span-5 lg:col-span-6 mb-8 md:mb-0 md:order-1"
              ref={tabs}
            >
              <div className="relative flex flex-col text-center lg:text-right">
                {features.map(({ id, name, imageUrl }) => (
                  <Transition
                    show={tab === id}
                    key={id}
                    appear={true}
                    className="w-full"
                    enter="transition ease-in-out duration-700 transform order-first"
                    enterStart="opacity-0 translate-y-16"
                    enterEnd="opacity-100 translate-y-0"
                    leave="transition ease-in-out duration-300 transform absolute"
                    leaveStart="opacity-100 translate-y-0"
                    leaveEnd="opacity-0 -translate-y-16"
                  >
                    <div className="relative inline-flex flex-col">
                      <img className="md:max-w-none mx-auto rounded" src={imageUrl} width="500" alt={name} />
                    </div>
                  </Transition>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

export default FeaturesHome;
