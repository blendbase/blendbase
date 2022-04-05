import { Html, Head, Main, NextScript } from "next/document";

export default function Document() {
  return (
    <Html>
      <Head>
        <meta charset="UTF-8" />
        <meta name="description" content="Single open-source GraphQL API to connect CRMs to your SaaS" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <link rel="icon" href="/favicon.png" />
      </Head>
      <body className="font-space-grotesk antialiased bg-white text-gray-900 tracking-tight">
        <Main />
        <NextScript />
      </body>
    </Html>
  );
}
