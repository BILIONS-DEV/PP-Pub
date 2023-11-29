const path = require('path');
const StringReplacePlugin = require("string-replace-webpack-plugin");
const UglifyJsPlugin = require('uglifyjs-webpack-plugin');
const HTMLWebpackPlugin = require('html-webpack-plugin');
const webpack = require('webpack');

const Config = (env, argv) => {
    env = env || {};
    console.log("===============");
    console.log("||Environment||");
    console.log("===============");
    console.log(env);
    console.log("===============");
    let htmlTemplate = "index";
    let entry = "thuy";
    let fileName = entry;
    let Plugins = [];
    let Rules = [];
    let Optimization = {};
    let Mode = argv.mode ? argv.mode : 'production';

    if (env.template) {
        htmlTemplate = env.template;
    }
    //build preview html file
    Plugins.push(new HTMLWebpackPlugin({
        template: path.resolve(__dirname, 'src/html/' + htmlTemplate + '.html')
    }));

    //babel loader
    Rules.push({
        test: /\.js$/,
        include: [
            path.resolve(__dirname, 'src'),
            path.resolve(__dirname, 'node_modules/observe-element-in-viewport/'),
            path.resolve(__dirname, 'node_modules/quadtree-lib'),
        ],
        exclude: /node_modules\/(?!(observe-element-in-viewport|quadtree-lib)).+/,
        use: {
            loader: "babel-loader",
            options: {
                presets: ["@babel/preset-env"]
            }
        }
    });

    //remove vilog
    if (env.nolog && env.nolog === "true") {
        Rules.push({
            test: /\.js$/,
            exclude: /node_modules/,
            loader: StringReplacePlugin.replace({
                replacements: [{
                    pattern: /vilog\(.*\);/g,
                    replacement: function (match, p1, offset, string) {
                        return "";
                    }
                }]
            })
        });
    }


    //entry
    if (env.entry) {
        entry = env.entry;
        fileName = entry;
    }

    //filename
    if (env.min && env.min === "true") {
        fileName += "_min";
    } else {
        fileName += "_max";
    }
    if (env.nolog && env.nolog === "true") {
        fileName += "_nolog";
    } else {
        fileName += "_log";
    }

    //minify
    if (env.min && env.min === "true") {
        Optimization.minimizer = [
            new UglifyJsPlugin()
        ];
    } else {
        Optimization.minimize = false;
    }

    let config = {};
    config.mode = Mode;
    config.entry = './src/' + entry + '.js';
    config.output = {
        path: path.resolve(__dirname, 'dist'),
        filename: fileName + '.js'
    };
    config.module = {
        rules: Rules
    };
    config.resolve = {
        extensions: ['.js']
    };
    config.devServer = {
        host: "jstag.local",
        port: 8080,
        contentBase: [__dirname + '/dist', __dirname + '/src/html'],
        inline: true,
        hot: true,
        open: false,
        disableHostCheck: true
    };
    config.plugins = Plugins;
    config.optimization = Optimization;
    return config;
};
module.exports = Config;