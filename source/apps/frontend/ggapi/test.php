<?php

ini_set("display_errors", 1);
$val = getopt(null, ["refreshToken:", "networkId:", "networkName:"]);
$networkName = '';
if (empty($tokenTest)) {
    if (empty($val["refreshToken"])) {
        echo json_encode(["status" => false, "message" => "refresh token is required"]);
        return;
    }
    if (empty($val["networkId"])) {
        echo json_encode(["status" => false, "message" => "networkId is required"]);
        return;
    }
    if ($val["refreshToken"]) {
        $refreshToken = $val["refreshToken"];
    }
    if ($val["networkId"]) {
        $networkId = $val['networkId'];
    }
    if (!empty($val["networkName"])) {
        if ($val["networkName"]) {
            $networkName = $val['networkName'];
        }
    }
} else {
    $refreshToken = $tokenTest;
    $networkId = '325081995';
    $networkName = 'X8';
}

if ($networkName == '') {
    $networkName = "PubPowerCheckAPI";
}

require __DIR__ . "/googleads-php-lib/vendor/autoload.php";

use Google\Auth\OAuth2;
use Google\AdsApi\AdManager\AdManagerSessionBuilder;
use Google\AdsApi\AdManager\v202302\ServiceFactory;

//https://myaccount.google.com/u/0/permissions?pli=1
$oauth2 = new OAuth2([
    'authorizationUri' => 'https://accounts.google.com/o/oauth2/v2/auth',
    'tokenCredentialUri' => 'https://www.googleapis.com/oauth2/v4/token',
    'redirectUri' => "https://apps.valueimpression.com/gam/callback",
    'clientId' => '1033050029778-nv67vpjb627a68e7kur22sn1s0so2v2s.apps.googleusercontent.com',
    'clientSecret' => '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',
    'scope' => 'https://www.googleapis.com/auth/dfp',
    // more
    'prompt' => 'consent',
    'approval_prompt' => 'force',
    'access_type' => 'offline'
]);

//$code = "4/0AX4XfWgex9AWXEPAuL2jV7n1nEcFfGlOtWvVs3DYrlpIJPgH7laFf6GUlN8JpkjTvf75aw";
//$oauth2->setCode($code);
// refresh token
//$access_token = "ya29.a0ARrdaM_24NHQxUZJws3i8E7BnOpb7U7MB0DDMQBzp-Cp0DNRi4NohZVGJ_sx79rUp3OXp6H1lcT8qxyLB7acHJOdPVBTssIKdP0VzMuXHwAA94eEPn-7g26y3v4_LNOu2zkOKiiAvyWIxtWvLygK5RlcMMg8";
//$refreshToken = "1//0dsYIgu1w4_5uCgYIARAAGA0SNwF-L9Ircv8f5MEY6EWYsMdStEgYscP6zwgTyOnQz6U08fIohTqzgNsQGfScSFWDn_a_072xk0E";
$oauth2->setGrantType("refresh_token");
$oauth2->setRefreshToken($refreshToken);
$oauth2->generateCredentialsRequest();

try {
    $authToken = $oauth2->fetchAuthToken();
    //     var_dump($authToken);
    //     $is = $oauth2->isExpired();
    //     var_dump($is);
} catch (Exception $exc) {
    echo json_encode(["status" => false, "type" => "refresh_token", "message" => $exc->getMessage()]);
    //     var_dump($exc->getMessage());
    //     var_dump($exc->getCode());
    //     var_dump($exc->getLine());
    return;
}

try {
    $session = (new AdManagerSessionBuilder())
        ->withOAuth2Credential($oauth2)
        ->withNetworkCode($networkId)
        ->withApplicationName($networkName)
        ->build();
} catch (Exception $exc) {
    echo json_encode(["status" => false, "type" => "oauth2", "message" => $exc->getMessage()]);
    return;
}

$serviceFactory = new ServiceFactory();
$networkService = $serviceFactory->createNetworkService($session);
try {
    $network = $networkService->getCurrentNetwork();
    var_dump($network);
} catch (\Google\AdsApi\AdManager\v202302\ApiException $e) {
    echo json_encode(["status" => false, "type" => "refresh_token", "message" => $e->getMessage()]);
    return;
}


echo json_encode(["status" => true, "message" => "Api Access Enabled!"]);
