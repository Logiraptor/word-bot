const { createConfig, defineConstants, env, entryPoint, setOutput, sourceMaps } = require("@webpack-blocks/webpack2");
const devServer = require("@webpack-blocks/dev-server2");
const typescript = require("@webpack-blocks/typescript");
const autoprefixer = require("autoprefixer");

module.exports = createConfig([
    entryPoint("./src/main.tsx"),
    setOutput("./public/bundle.js"),
    typescript(),
    defineConstants({
        "process.env.NODE_ENV": process.env.NODE_ENV,
    }),
    env("development", [
        devServer(),
        devServer.proxy({
            "/api": { target: "http://localhost:3000" },
        }),
        sourceMaps(),
    ]),
]);
