/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  async redirects() {
    return [
      {
        source: "/",
        destination: "/integrations",
        permanent: false
      }
    ];
  }
};

module.exports = nextConfig;
