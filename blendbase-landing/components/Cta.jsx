import React from "react";
import AppConfig from "../app.config";

function Cta() {
  return (
    <section className="mt-24">
      <div className="max-w-6xl mx-auto px-4 sm:px-6">
        <div className="pb-12 md:pb-20">
          {/* CTA box */}
          <div className="bg-blue-600 rounded py-10 px-8 md:py-16 md:px-12 shadow-2xl">
            <div className="flex flex-col lg:flex-row justify-between items-center">
              {/* CTA content */}
              <div className="mb-6 lg:mr-16 lg:mb-0 text-center lg:text-left">
                <h3 className="h3 text-white mb-2">Ready to get started?</h3>
                <p className="text-white text-lg opacity-75">
                  Go to Github, clone the project and join our Discord to hit the ground running
                </p>
              </div>

              {/* CTA button */}
              <div>
                <a
                  className="btn text-blue-600 hover:bg-white bg-brand-50 hover:from-blue-100 transition-all duration-200"
                  href={AppConfig.githubUrl}
                >
                  Get started
                </a>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

export default Cta;
