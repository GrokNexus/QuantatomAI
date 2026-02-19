/** @type {import('next').NextConfig} */
const nextConfig = {
    reactStrictMode: true,
    // Layer 6.2: Enable WebGPU (if needed via flags, though usually standard now)
    // webpack: (config) => {
    //   return config;
    // },
};

module.exports = nextConfig;
