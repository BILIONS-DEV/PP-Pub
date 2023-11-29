<?php
session_start();

/************************************************
ATTENTION: Change this path to point to your vendor folder if your project
directory structure differs from this repository's!
 ************************************************/
require_once __DIR__ . '/google-ads-api/vendor/autoload.php';

// Autoload example classes.
spl_autoload_register(function ($class_name) {
    include 'examples/' . $class_name . '.php';
});

// Max results per page.
define('MAX_LIST_PAGE_SIZE', 50, true);
define('MAX_REPORT_PAGE_SIZE', 50, true);

// Configure token storage on disk.
// If you want to store refresh tokens in a local disk file, set this to true.
define('STORE_ON_DISK', false, true);
define('TOKEN_FILENAME', 'tokens.dat', true);
$refreshToken = "1//0fyTOoXTKmMXRCgYIARAAGA8SNwF-L9Irlh6qCUxexGeDuF1HveXkf9p93JQBqXVTHQU2pBGRE7VfHs-8TIjbl0tyvyagGiqoYuo";
$client = new Google_Client();
$client->addScope('https://www.googleapis.com/auth/adsense.readonly');
$client->setAccessType('offline');
$client->setApprovalPrompt('force');
try {
    $client->setAuthConfig("client_secret_1032228041816-u7cg9gr6n0bubs8i5iegv4j858m2506v.apps.googleusercontent.com.json");
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

$accounts = getAllAccounts($service, MAX_LIST_PAGE_SIZE);
var_dump($accounts);
if (isset($accounts) && !empty($accounts)) {
    // Get an example account, so we can run the following sample.
    $exampleAccountId = $accounts[0]['name'];
    GetAccountTree::run($service, $exampleAccountId);
    $adClients =
        GetAllAdClients::run($service, $exampleAccountId, MAX_LIST_PAGE_SIZE);

    if (isset($adClients) && !empty($adClients)) {
        // Get an ad client ID, so we can run the rest of the samples.
        $exampleAdClient = end($adClients);
        $exampleAdClientId = $exampleAdClient['name'];

        GenerateReport::run($service, $exampleAccountId, $exampleAdClientId);
    } else {
        print "No ad clients found, unable to run dependent examples.\n";
    }

} else {
    print 'No accounts found, unable to run dependant examples.\n';
}

function getAllAccounts($service, $maxPageSize) {
    $optParams['pageSize'] = $maxPageSize;

    $pageToken = null;
    do {
        $optParams['pageToken'] = $pageToken;
        $result = $service->accounts->listAccounts($optParams);
        $accounts = null;
        if (!empty($result['accounts'])) {
            $accounts = $result['accounts'];
            $pageToken = $result['nextPageToken'];
        }
    } while ($pageToken);

    return $accounts;
}