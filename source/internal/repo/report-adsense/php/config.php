<?php

$val = getopt(null, ["refreshToken:","accountAdsense:"]);

if (empty($val["refreshToken"])) {
    echo json_encode(["status" => false, "message" => "refresh token is required"]);
    return;
}
if ($val["refreshToken"]) {
    $refreshToken = $val["refreshToken"];
}

if (empty($val["accountAdsense"])) {
    echo json_encode(["status" => false, "message" => "accountAdsense is required"]);
    return;
}

if ($val["accountAdsense"]) {
    $accountAdsense = $val["accountAdsense"];
}

session_start();

/************************************************
ATTENTION: Change this path to point to your vendor folder if your project
directory structure differs from this repository's!
 ************************************************/
require_once __DIR__ . '/google-ads-api/vendor/autoload.php';

// Max results per page.

// Configure token storage on disk.
// If you want to store refresh tokens in a local disk file, set this to true.
$client = new Google_Client();
$client->addScope('https://www.googleapis.com/auth/adsense.readonly');
$client->setAccessType('offline');
$client->setApprovalPrompt('force');
try {
    $client->setAuthConfig(__DIR__."/app_".$accountAdsense.".json");
} catch (\Google\Exception $e) {
    var_dump($e);
}
$service = new Google_Service_Adsense($client);
$client->refreshToken($refreshToken);

//echo "EXPIRE IN " . ($token["expires_in"]) . PHP_EOL;
if ($client->isAccessTokenExpired()) {
//    echo "TOKEN EXPIRED" . PHP_EOL;
    $client->fetchAccessTokenWithRefreshToken($refreshToken);
    $client->setAccessToken($client->getAccessToken());
    $newToken = $client->getAccessToken();
//    echo "NEW TOKEN SAVED" . PHP_EOL;
//    echo "EXPIRE IN " . ($newToken["expires_in"]) . PHP_EOL;
}