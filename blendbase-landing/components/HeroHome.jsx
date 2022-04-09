import App from "next/app";
import React, { useState } from "react";
import AppConfig from "../app.config";

function HeroHome() {
  const [videoModalOpen, setVideoModalOpen] = useState(false);

  return (
    <section className="relative">
      <div className="absolute w-full  top-16 pointer-events-none -z-1" aria-hidden="true">
        <img src="/images/landing-hero-back.png" alt="Hero background" width="100%" />
      </div>

      <div className="max-w-6xl mx-auto px-4 sm:px-6">
        {/* Hero content */}
        <div className="pt-32 md:pt-56">
          {/* Section header */}
          <div className="text-center pb-12 md:pb-32">
            <h1 className="text-5xl font-extrabold leading-none tracking-tighter">
              Build SaaS integrations <br />
              <span className="bg-clip-text text-transparent bg-gradient-to-r from-brand-500 to-pink-500">
                at scale
              </span>
            </h1>
            <div className="max-w-3xl mx-auto mt-8">
              <p className="text-xl text-gray-600 mb-8">Single open-source GraphQL API to connect CRMs to your SaaS</p>
              <div className="max-w-xs mx-auto sm:max-w-none sm:flex sm:justify-center">
                <div>
                  <a
                    className="btn btn-lg text-white bg-brand-500 hover:bg-brand-700 w-full mb-4 sm:w-auto sm:mb-0"
                    href={AppConfig.documentationUrl}
                    target="_blank"
                  >
                    Get started
                  </a>
                </div>
                <div>
                  <a
                    className="btn btn-lg text-white bg-slate-700 hover:bg-slate-900 w-full sm:w-auto sm:ml-4"
                    target="_blank"
                    href={AppConfig.discordUrl}
                  >
                    Join our Discord
                  </a>
                </div>
              </div>
            </div>
          </div>

          {/* Hero image */}
          <div className="mt-16  rounded rounded-2xl py-8 px-12">
            <div className="relative flex justify-center">
              <div className="flex flex-col justify-center">
                <img className="mx-auto w-full" src="/images/hero-schema.png" alt="Hero" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

export default HeroHome;
