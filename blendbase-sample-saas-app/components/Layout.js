/* This example requires Tailwind CSS v2.0+ */
import { Fragment, useState } from "react";
import { Dialog, Transition } from "@headlessui/react";
import { CogIcon, UserGroupIcon } from "@heroicons/react/outline";

import classNames from "../utils/classNames";

const crmNavigation = [{ name: "Contacts", href: "/contacts", icon: UserGroupIcon }];

const settingsNavigation = [{ name: "Integrations", href: "/integrations", icon: CogIcon }];

function MenuItem({ name, href, icon: Icon }) {
  return (
    <a
      key={name}
      href={href}
      className={classNames(
        "text-gray-600 hover:bg-gray-50 hover:text-gray-900",
        "group flex items-center rounded-md px-2 py-2 text-sm font-medium"
      )}
    >
      <Icon
        className={classNames("text-gray-400 group-hover:text-gray-500", "mr-3 h-6 w-6 flex-shrink-0")}
        aria-hidden="true"
      />
      {name}
    </a>
  );
}

export default function Layout({ children }) {
  return (
    <>
      <div>
        <div className="hidden md:fixed md:inset-y-0 md:flex md:w-64 md:flex-col">
          <div className="flex min-h-0 flex-1 flex-col border-r border-gray-200 bg-white">
            <div className="flex flex-1 flex-col overflow-y-auto pt-5 pb-4">
              <div className="flex flex-shrink-0 items-center px-4">
                <img
                  className="h-8 w-auto"
                  src="https://tailwindui.com/img/logos/workflow-logo-indigo-600-mark-gray-800-text.svg"
                  alt="Workflow"
                />
              </div>
              <nav className="mt-12 bg-white px-2">
                <h3
                  className="px-3 text-xs font-semibold uppercase tracking-wider text-gray-500"
                  id="projects-headline"
                >
                  CRM
                </h3>
                <div className="flex-1 space-y-1">
                  {crmNavigation.map((item) => (
                    <MenuItem key={item.name} {...item} />
                  ))}
                </div>
                <div className="mt-8">
                  <h3
                    className="px-3 text-xs font-semibold uppercase tracking-wider text-gray-500"
                    id="projects-headline"
                  >
                    Settings
                  </h3>
                  <div className="mt-1 space-y-1" aria-labelledby="projects-headline">
                    {settingsNavigation.map((item) => (
                      <MenuItem key={item.name} {...item} />
                    ))}
                  </div>
                </div>
              </nav>

              <nav className="mt-5 flex-1 space-y-1 bg-white px-2"></nav>
            </div>
          </div>
        </div>
        <div className="flex flex-1 flex-col md:pl-64">
          <main className="flex-1">
            <div className="p-8">{children}</div>
          </main>
        </div>
      </div>
    </>
  );
}
