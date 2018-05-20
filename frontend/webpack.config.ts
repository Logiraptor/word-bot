import * as ExtractTextPlugin from "extract-text-webpack-plugin";
import * as HtmlWebpackPlugin from "html-webpack-plugin";
import * as webpack from "webpack";

declare const __dirname: string;

const config: webpack.Configuration = {
    entry: __dirname + "/main.tsx",
    output: {
        path: __dirname + "/public",
        filename: "app.bundle.js",
    },
    resolve: {
        extensions: [ ".scss", ".ts", ".tsx", ".js", ".go" ],
    },
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                use: "awesome-typescript-loader",
            },
            {
                test: /\.go$/,
                use: "gopherjs-loader",
            },
            {
                test: /\.scss$/,
                use: ExtractTextPlugin.extract({
                    fallback: "style-loader",
                    use: [ "css-loader", "sass-loader" ],
                }),
            },
        ],
    },
    plugins: [ new ExtractTextPlugin("style.css"), new HtmlWebpackPlugin() ],
};

export default config;
