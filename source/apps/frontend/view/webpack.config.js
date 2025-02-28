const path = require('path');
const TerserPlugin = require("terser-webpack-plugin");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const CssMinimizerPlugin = require("css-minimizer-webpack-plugin");

const Config = (env, argv) => {
    console.log(argv)
    console.log(`This is the Webpack 5 'mode': ${argv.mode} for frontend`);
    let mode = argv.mode ? argv.mode : 'production';
    let configObject = {
        entry: {
            // "css/style": ["./pages/_css/style.css"],
            // Test
            // "home/test.min": ["./pages/home/test.js", "./pages/home/test.ts"],
            // "home/test.min": ["./pages/home/test.js"],
            // Libs
            // "js/functions.min": ["./pages/js_libs/functions.js"],
            // "be/supply/index": ["./pages/_backend/supply/index.js"],

            //=> Page:  AdsTxt
            // "ads_txt/index": ["./pages/ads_txt/index.js", "./pages/ads_txt/index.css"],
            // "ads_txt/detail": ["./pages/ads_txt/detail.js"],
            //=> Page: AdTag
            // "adtag/add": ["./pages/adtag/add.js"],
            // "adtag/edit": ["./pages/adtag/edit.js"],
            //=> Page: Bidder
            // "bidder/index": ["./pages/bidder/index.js"],
            // "bidder/add": ["./pages/bidder/add.js"],
            // "bidder/edit": ["./pages/bidder/edit.js"],
            // "bidder/view": ["./pages/bidder/view.js"],
            //=> Page: Blocking
            // "blocking/index": ["./pages/blocking/index.js"],
            // "blocking/add": ["./pages/blocking/add.js"],
            // "blocking/edit": ["./pages/blocking/edit.js"],
            //=> Page: Config
            "config/index": ["./pages/config/index.js"],
            //=> Page: Content
            // "content/index": ["./pages/content/index.js"],
            // "content/add": ["./pages/content/add.js"],
            // "content/add/quiz": ["./pages/content/add.js"],
            // "content/add/video": ["./pages/content/add.js"],
            // "content/edit": ["./pages/content/edit.js"],
            // "content/edit/quiz": ["./pages/content/edit.js"],
            // "content/edit/video": ["./pages/content/edit.js"],
            //=> Page: Floor
            // "floor/index": ["./pages/floor/index.js"],
            // "floor/add": ["./pages/floor/add.js"],
            // "floor/edit": ["./pages/floor/edit.js"],
            //=> Page: GAM
            // "gam/index": ["./pages/gam/index.js"],
            // "gam/edit": ["./pages/gam/edit.js"],
            //=> Page: Home
            "home/dashboard": ["./pages/home/dashboard.js"],
            //=> Page: Inventory
            "websites/index": ["./pages/websites/index.js"],
            "websites/setup": ["./pages/websites/setup.js"],
            "websites-v2/setup": ["./pages/websites_v2/setup.js"],
            //=> Page: LineItem
            // "line-item/index": ["./pages/line-item/index.js"],
            // "line-item/add": ["./pages/line-item/add.js"],
            // "line-item/edit": ["./pages/line-item/edit.js"],
            //=> Page: LineItem
            // "line-item-v2/index": ["./pages/line-item-v2/index.js"],
            // "line-item-v2/add": ["./pages/line-item-v2/add.js"],
            // "line-item-v2/edit": ["./pages/line-item-v2/edit.js"],
            //=> Page: Player/Template
            // "player/index": ["./pages/player/index.js"],
            // "player/add": ["./pages/player/add.js"],
            // "player/edit": ["./pages/player/edit.js"],
            // "player/view": ["./pages/player/view.js"],

            // "player-v2/index": ["./pages/player-v2/index.js"],
            // "player-v2/add": ["./pages/player-v2/add.js"],
            // "player-v2/edit": ["./pages/player-v2/edit.js"],
            // "player-v2/view": ["./pages/player-v2/view.js"],

            //=> Page: Playlist
            // "playlist/index": ["./pages/playlist/index.js"],
            // "playlist/add": ["./pages/playlist/add.js"],
            // "playlist/edit": ["./pages/playlist/edit.js"],
            // "playlist/view": ["./pages/playlist/view.js"],
            //=> Page: User
            "user/account": ["./pages/user/account.js"],
            // "user/billing": ["./pages/user/billing.js"],
            "user/forget-password": ["./pages/user/forget-password.js"],
            "user/login": ["./pages/user/login.js"],
            "user/register": ["./pages/user/register.js"],
            // "user/password": ["./pages/user/password.js"],
            "user/new-password": ["./pages/user/new-password.js"],
            //=> Page: Video
            // "video/index": ["./pages/video/index.js"],
            //=> Page: Video
            // "support/index": ["./pages/support/index.css", "./pages/support/index.js"],
            // "support/product-description": ["./pages/support/index.css"],
            // "support/tickets/new": ["./pages/support/index.css", "./pages/support/index.js"],
            //=> Page: Identity
            // "identity/index": ["./pages/identity/index.js"],
            // "identity/add": ["./pages/identity/add.js"],
            // "identity/edit": ["./pages/identity/edit.js"],
            //=> Page: Channels
            // "channels/index": ["./pages/channels/index.js"],
            // "channels/add": ["./pages/channels/add.js"],
            // "channels/edit": ["./pages/channels/edit.js"],
            //=> Page: A/B Testing
            // "ab_testing/index": ["./pages/ab_testing/index.js"],
            // "ab_testing/add": ["./pages/ab_testing/add.js"],
            // "ab_testing/edit": ["./pages/ab_testing/edit.js"],
            //=> Page: Payment
            "payment": ["./pages/payment/index.js"],
            "payment/index": ["./pages/payment/index.js"],
            //=> Page: Rule
            // "rule/index": ["./pages/rule/index.js"],
            //=> Page: Rule
            // "blocked-page/add": ["./pages/blocked-page/add.js"],
            // "blocked-page/edit": ["./pages/blocked-page/edit.js"],
            //=> Page: Activity
            "history/index": ["./pages/history/index.js"],
            //=> Page: Content Quiz
            // "quiz/index": ["./pages/quiz/index.js"],
            // "content/add-quiz": ["./pages/content/add-quiz.js"],
            // "content/edit-quiz": ["./pages/content/edit-quiz.js"],
            //=> Page: Ad Schedules
            // "ad_schedules/index": ["./pages/ad_schedules/index.js"],
            // "ad_schedules/add": ["./pages/ad_schedules/add.js"],
            // "ad_schedules/edit": ["./pages/ad_schedules/edit.js"],
            //=> Page: Campaign
            // "campaigns": ["./pages/campaign/index.js"],
            // "campaigns/index": ["./pages/campaign/index.js"],
            //=> Page: Ad Block
            // "adblock/analytics": ["./pages/adblock/analytics.js"],
            // "adblock/generator": ["./pages/adblock/generator.js"],
            //=> Page: Report
            // "report/index": ["./pages/report/analytics.js"],
            // "report/dimension": ["./pages/report/generator.js"],
            // "report/saved": ["./pages/report/generator.js"],
        },
        output: {
            // filename: '[name].min.js',
            filename: '[name].js',
            path: path.resolve(__dirname, '.assets')
        },
        plugins: [
            new MiniCssExtractPlugin({
                // filename: '[name].min.css'
                filename: '[name].css'
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
                            presets: ["@babel/preset-env", "@babel/preset-react"],
                            plugins: [
                                "@babel/plugin-transform-class-properties",
                                "@babel/plugin-proposal-object-rest-spread"
                            ]
                        }
                    }
                },
                {
                    test: /\.tsx?/,
                    use: {
                        loader: 'ts-loader',
                        options: {
                            transpileOnly: true,
                        }
                    },
                    exclude: /node_modules/,
                },
                {
                    // test: /\.css$/i,
                    test: /.s?css$/,
                    use: [
                        mode === 'production' ? MiniCssExtractPlugin.loader : 'style-loader',
                        // MiniCssExtractPlugin.loader, //=> sử dụng cái này để minify css vào file riêng
                        // "style-loader", //=> sử dụng cái này nếu muốn nhúng CSS vào JS
                        "css-loader",
                    ],
                },
            ]
        },
        resolve: {
            extensions: ['.tsx', '.ts', '.js'],
        },
    }
    let Optimization = {};
    if (mode === "production") {
        Optimization.minimize = true
        Optimization.minimizer = [
            new TerserPlugin(),
            new CssMinimizerPlugin()
        ];
    } else {
        Optimization.minimize = false;
    }
    configObject.optimization = Optimization
    // configObject.entry = {
    //     "test/index": ["./pages/test/index.ts", "./pages/test/index.css"],
    // }
    return configObject
}
module.exports = Config;