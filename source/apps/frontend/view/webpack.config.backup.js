const path = require('path');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const OptimizeCSSAssetsPlugin = require('optimize-css-assets-webpack-plugin');
const TerserPlugin = require("terser-webpack-plugin");

module.exports = (_, {mode}) => ({
    entry: {
        // Test
        "home/test.min": ["./pages/home/test.js"],
        // Libs
        // "js/functions.min": ["./pages/js_libs/functions.js"],

        //=> Page:  AdsTxt
        "ads_txt/index": ["./pages/ads_txt/index.js", "./pages/ads_txt/index.css"],
        "ads_txt/detail": ["./pages/ads_txt/detail.js"],
        //=> Page: AdTag
        "adtag/add": ["./pages/adtag/add.js"],
        "adtag/edit": ["./pages/adtag/edit.js"],
        //=> Page: Bidder
        "bidder/index": ["./pages/bidder/index.js"],
        "bidder/add": ["./pages/bidder/add.js"],
        "bidder/edit": ["./pages/bidder/edit.js"],
        //=> Page: Blocking
        "blocking/index": ["./pages/blocking/index.js"],
        "blocking/add": ["./pages/blocking/add.js"],
        "blocking/edit": ["./pages/blocking/edit.js"],
        //=> Page: Config
        "config/index": ["./pages/config/index.js"],
        //=> Page: Content
        "content/index": ["./pages/content/index.js"],
        "content/add": ["./pages/content/add.js"],
        "content/edit": ["./pages/content/edit.js"],
        //=> Page: Floor
        "floor/index": ["./pages/floor/index.js"],
        "floor/add": ["./pages/floor/add.js"],
        "floor/edit": ["./pages/floor/edit.js"],
        //=> Page: GAM
        "gam/index": ["./pages/gam/index.js"],
        "gam/edit": ["./pages/gam/edit.js"],
        //=> Page: Home
        "home/dashboard": ["./pages/home/dashboard.js"],
        //=> Page: Inventory
        "supply/index": ["./pages/supply/index.js"],
        "supply/setup": ["./pages/supply/setup.js"],
        //=> Page: LineItem
        "line-item/index": ["./pages/line-item/index.js"],
        "line-item/add": ["./pages/line-item/add.js"],
        "line-item/edit": ["./pages/line-item/edit.js"],
        //=> Page: Player/Template
        "player/index": ["./pages/player/index.js"],
        "player/add": ["./pages/player/add.js"],
        "player/edit": ["./pages/player/edit.js"],
        "player/view": ["./pages/player/view.js"],
        //=> Page: Playlist
        "playlist/index": ["./pages/playlist/index.js"],
        "playlist/add": ["./pages/playlist/add.js"],
        "playlist/edit": ["./pages/playlist/edit.js"],
        //=> Page: User
        "user/account": ["./pages/user/account.js"],
        "user/billing": ["./pages/user/billing.js"],
        "user/forget-password": ["./pages/user/forget-password.js"],
        "user/login": ["./pages/user/login.js"],
        "user/register": ["./pages/user/register.js"],
        "user/password": ["./pages/user/password.js"],
        "user/new-password": ["./pages/user/new-password.js"],
        //=> Page: Video
        "video/index": ["./pages/video/index.js"],
    },
    optimization: {
        minimizer: [
            new TerserPlugin(),
            new OptimizeCSSAssetsPlugin({})
        ],
    },
    output: {
        // filename: '[name].min.js',
        filename: '[name].js',
        path: path.resolve(__dirname, '.assets')
    },
    plugins: [
        new MiniCssExtractPlugin({
            filename: '[name].min.css'
        })
    ],
    module: {
        rules: [
            {
                test: /\.m?js$/,
                exclude: /(node_modules|bower_components)/,
                use: {
                    loader: 'babel-loader',
                    options: {
                        presets: ['@babel/preset-env'],
                        plugins: ['@babel/plugin-proposal-object-rest-spread']
                    }
                }
            },
            {
                test: /\.css$/i,
                use: [
                    mode === 'production' ? MiniCssExtractPlugin.loader : 'style-loader',
                    // MiniCssExtractPlugin.loader, //=> sử dụng cái này để minify css vào file riêng
                    // "style-loader", //=> sử dụng cái này nếu muốn nhúng CSS vào JS
                    "css-loader",
                ],
            },
        ]
    }
});

// module.exports = {
//     entry: {
//         "pomodoro/index": ["./Pomodoro/index.js", "./Pomodoro/index.css"],
//         "pomodoro/test": "./Pomodoro/test.js",
//     },
//     optimization: {
//         minimizer: [
//             new TerserPlugin(),
//             new OptimizeCSSAssetsPlugin({})
//         ],
//     },
//     output: {
//         filename: '[name].min.js',
//         path: path.resolve(__dirname, '.outputs')
//     },
//     plugins: [
//         new MiniCssExtractPlugin({
//             filename: '[name].min.css'
//         })
//     ],
//     module: {
//         rules: [
//             {
//                 test: /\.m?js$/,
//                 exclude: /(node_modules|bower_components)/,
//                 use: {
//                     loader: 'babel-loader',
//                     options: {
//                         presets: ['@babel/preset-env'],
//                         plugins: ['@babel/plugin-proposal-object-rest-spread']
//                     }
//                 }
//             },
//             {
//                 test: /\.css$/i,
//                 use: [
//                     // MiniCssExtractPlugin.loader,
//                     "style-loader", //=> sử dụng cái này nếu muốn nhúng CSS vào JS
//                     "css-loader",
//                 ],
//             },
//         ]
//     }
// };