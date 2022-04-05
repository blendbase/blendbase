import Head from "next/head";

import Header from "../components/Header";
import HeroHome from "../components/HeroHome";
import Workflows from "../components/Workflows";
import FeaturesHome from "../components/FeaturesHome";
import Footer from "../components/Footer";
import Cta from "../components/Cta";

export default function Home() {
  return (
    <div className="flex flex-col min-h-screen overflow-hidden">
      <Head>
        <title>Blendbase – Single CRM API</title>
      </Head>
      <Header />

      {/*  Page content */}
      <main className="grow">
        {/*  Page sections */}
        <HeroHome />
        <Workflows />
        <FeaturesHome />
        <Cta />
      </main>

      {/*  Site footer */}
      <Footer />
    </div>
  );
}
