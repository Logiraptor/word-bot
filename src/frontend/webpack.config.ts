import * as webpack from "webpack";
import * as HtmlWebpackPlugin from "html-webpack-plugin";
import * as ExtractTextPlugin from "extract-text-webpack-plugin";

let __dirname: string;

const config: webpack.Configuration = {
    entry: __dirname + "/main.tsx",
    output: {
        path: __dirname + "/public",
        filename: "app.bundle.js",
    },
    resolve: {
        extensions: [ ".scss", ".ts", ".tsx", ".js" ],
    },
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                use: "awesome-typescript-loader",
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
    plugins: [ new HtmlWebpackPlugin(), new ExtractTextPlugin("style.css") ],
};

export default config;
