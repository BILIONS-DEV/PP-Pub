#!/bin/bash

#Define constants
LOG_PREFIX="=======[LOG]======="
#TEMP_PREBID_PATH=Prebid.js_Temp
TEMP_PREBID_PATH=Prebid.js_Temp

# exit when any command fails
set -e
# keep track of the last executed command
trap 'last_command=$current_command; current_command=$BASH_COMMAND' DEBUG
# echo an error message before exiting
trap 'echo "$LOG_PREFIX \"${last_command}\" command filed with exit code $?."' EXIT

#Clean temp directory
echo "$LOG_PREFIX Clean temp directory"
rm -rf $TEMP_PREBID_PATH
sleep 1

#Clone latest Prebid.js
LATEST_RELEASE=$(curl -L -s -H 'Accept: application/json' https://github.com/prebid/Prebid.js/releases/latest)
LATEST_VERSION=$(echo $LATEST_RELEASE | sed -e 's/.*"tag_name":"\([^"]*\)".*/\1/')
echo "$LOG_PREFIX Clone Prebid.js latest version $LATEST_VERSION"
curl -L https://github.com/prebid/Prebid.js/archive/refs/tags/$LATEST_VERSION.zip --output prebid.zip
unzip -q prebid.zip
yes | mv Prebid.js-$LATEST_VERSION $TEMP_PREBID_PATH
yes | rm prebid.zip
cd $TEMP_PREBID_PATH

#Install package dependencies
echo "$LOG_PREFIX Install package dependencies"
npm install

#Replace globalVarName to apdpbjs
echo "$LOG_PREFIX Replace globalVarName to apdpbjs"
sed -i "s/\"globalVarName\": \"pbjs\"/\"globalVarName\": \"apdpbjs\"/g" package.json

#Add bidtestBidAdapter.js
echo "$LOG_PREFIX Add bidtestBidAdapter.js"
cat > modules/bidtestBidAdapter.js << EOF
import { registerBidder } from '../src/adapters/bidderFactory.js';
const BIDDER_CODE = 'bidtest';

const ENDPOINT = window._PBCFG.ENDPOINT;
const USER_SYNC_URL = window._PBCFG.USER_SYNC_URL;
const BRANDNAME = window._PBCFG.BRANDNAME;
const HOMEPAGE = window._PBCFG.HOMEPAGE;
const BANNER_BG = window._PBCFG.BANNER_BG;
const NATIVE_ICON = window._PBCFG.NATIVE_ICON;
const NATIVE_IMAGE = window._PBCFG.NATIVE_IMAGE;
const NATIVE_DISPLAY_URL = window._PBCFG.NATIVE_DISPLAY_URL;
const NATIVE_PRIVACY_ICON = window._PBCFG.NATIVE_PRIVACY_ICON;

var videoLibrary = [
  'https://s0.2mdn.net/4253510/google_ddm_animation_480P.mp4',
  'https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4',
  'https://s-static.innovid.com/media/encoded/05_16/18049/7_source_27856_45954.mp4',
  'https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerBlazes.mp4',
  'https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4',
  'https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerFun.mp4',
  'https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerJoyrides.mp4',
  'https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerMeltdowns.mp4',
  'https://vcdn.adnxs.com/p/creative-video/c7/84/47/f8/c78447f8-84f7-400b-aacf-39a87d9e1fb7/c78447f8-84f7-400b-aacf-39a87d9e1fb7_768_432_500k.mp4'
];
function getWidth(sizes) {
  return sizes[0][0];
}
function getHeight(sizes) {
  return sizes[0][1];
}
export const spec = {
  code: BIDDER_CODE,
  supportedMediaTypes: ['banner', 'video', 'native'],
  isBidRequestValid: function (bid) {
    return true;
  },

  buildRequests: function (validBidRequests, bidderRequest) {
    const payload = bidderRequest;
    var data = {
      page: document.location.href,
      ref: document.referrer,
      domain: document.location.hostname,
      ua: navigator.userAgent,
      dnt: (navigator.doNotTrack === 'yes' || navigator.doNotTrack === '1' || navigator.msDoNotTrack === '1') ? 1 : 0,
      device_height: window.innerHeight,
      device_width: window.innerWidth,
      language: navigator.language
    };
    payload.data = data;
    const payloadString = JSON.stringify(payload);
    return {
      method: 'POST',
      url: ENDPOINT,
      data: payloadString
    };
  },
  interpretResponse: function (serverResponse, bidRequest) {
    const bidResponses = [];

    var request = JSON.parse(bidRequest.data)
    request.bids.forEach(bid => {
      var bidResponse = {};
      if (bid.mediaTypes.video) {
        bidResponse = {
          requestId: bid.bidId,
          cpm: 0.5,
          width: getWidth(bid.sizes),
          height: getHeight(bid.sizes),
          creativeId: 'video-preview',
          dealId: 'sample-deal-id',
          currency: 'USD',
          netRevenue: false,
          ttl: 360,
          mediaType: 'video',
          vastXml: '<?xml version="1.0" encoding="utf-8"?><VAST version="2.0"><Ad id="12345"><InLine><AdSystem version="1.0"><![CDATA[' + BRANDNAME + ']]></AdSystem><AdTitle><![CDATA[Sample VAST]]></AdTitle><Description><![CDATA[A sample VAST feed]]></Description><Creatives><Creative sequence="1" id="1"><Linear skipoffset="00:00:05"><Duration>00:00:30</Duration><TrackingEvents></TrackingEvents><VideoClicks><ClickThrough><![CDATA[' + HOMEPAGE + ']]></ClickThrough></VideoClicks><MediaFiles><MediaFile delivery="progressive" bitrate="256" width="640" height="480" type="video/mp4"><![CDATA[' + videoLibrary[Math.floor(Math.random() * videoLibrary.length)] + ']]></MediaFile></MediaFiles></Linear></Creative></Creatives></InLine></Ad></VAST>'
        };
      }
      if (bid.mediaTypes.banner) {
        bidResponse = {
          requestId: bid.bidId,
          cpm: 0.1,
          width: getWidth(bid.sizes),
          height: getHeight(bid.sizes),
          creativeId: 'banner-preview',
          dealId: 'sample-deal-id',
          currency: 'USD',
          netRevenue: false,
          ttl: 360,
          ad: '<div style="width: 100%; height: 100vh; background-color: ' + BANNER_BG + '; position: relative;"><div style="position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); text-align: center;">' + BRANDNAME + ' Preview Ad</div></div>'
        };
      }

      if (bid.mediaTypes.native) {
        bidResponse = {
          requestId: bid.bidId,
          cpm: 0.1,
          width: getWidth(bid.sizes),
          height: getHeight(bid.sizes),
          creativeId: 'native-preview',
          dealId: 'sample-deal-id',
          currency: 'USD',
          netRevenue: false,
          mediaType: 'native',
          ttl: 360,
          native: {
            title: BRANDNAME + ' native preview ad',
            body: 'Two things are infinite: the universe and human stupidity; and I\'m not sure about the universe.',
            body2: 'If you can\'t explain it to a six year old, you don\'t understand it yourself.',
            sponsoredBy: BRANDNAME,
            icon: NATIVE_ICON,
            image: NATIVE_IMAGE,
            clickUrl: HOMEPAGE,
            displayUrl: NATIVE_DISPLAY_URL,
            privacyLink: HOMEPAGE + '/privacy.html',
            privacyIcon: NATIVE_PRIVACY_ICON,
            cta: 'Learn More',
            rating: '5',
            downloads: '1000000',
            likes: '1000000',
            price: '1500',
            salePrice: '999',
            address: 'Mountain View, California, USA',
            phone: '123 456 789',
            clickTrackers: [
              HOMEPAGE
            ],
            impressionTrackers: [
              HOMEPAGE
            ],
            javascriptTrackers: '<script type="text/javascript" async="true" src="' + HOMEPAGE + '"></script>'
          }
        };
      }

      bidResponses.push(bidResponse);
    });
    return bidResponses;
  },
  getUserSyncs: function (syncOptions, serverResponses) {
    const syncs = [];
    if (syncOptions.iframeEnabled) {
      syncs.push({
        type: 'iframe',
        url: USER_SYNC_URL
      });
    }
    return syncs;
  },

  onTimeout: function (timeoutData) { },

  onBidWon: function (bid) { },

  onSetTargeting: function (bid) { }
};
registerBidder(spec);
EOF

#Add custom gulp task
# echo "$LOG_PREFIX Add custom gulp task"
# echo "gulp.task(makeWebpackPkg);" >> gulpfile.js
# echo "gulp.task('pbBundle', gulp.series(gulpBundle.bind(null, false)));" >> gulpfile.js

#Fix onemobile adapter
# echo "$LOG_PREFIX Fix onemobile adapter"
# sed -i "s/bidderCode === AOL_BIDDERS_CODES.ONEMOBILE/bidderCode === AOL_BIDDERS_CODES.ONEMOBILE || bidderCode === 'pp_onemobile'/g" modules/aolBidAdapter.js

#Fix sharethrough adapter
echo "$LOG_PREFIX Fix sharethrough adapter"
sed -i "s/\&\& bid.bidder === BIDDER_CODE/\&\& (bid.bidder === BIDDER_CODE || bid.bidder === 'pp_sharethrough')/g" modules/sharethroughBidAdapter.js

#Change Prebid log prefix
echo "$LOG_PREFIX Change Prebid log prefix"
sed -i "s/3b88c3/ffc000/g" src/utils.js
# sed -i "s/%cPrebid/%cPubPower Prebid/g" src/utils.js


#Compile Prebid.js
echo "$LOG_PREFIX Compile Prebid.js"
./node_modules/.bin/gulp build

cd -

#Remove old Prebid.js version
echo "$LOG_PREFIX Remove old Prebid.js version"
rm -rf Prebid.js
sleep 1

#Use new Prebid.js version
echo "$LOG_PREFIX Use new Prebid.js version"
mv $TEMP_PREBID_PATH Prebid.js
sleep 1

cd Prebid.js
echo "PREBID.JS PATH IS:"
pwd