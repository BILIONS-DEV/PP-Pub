<?php

ini_set("display_errors", 1);
// $tokenTest = "1//01jBGOVmHgt9gCgYIARAAGAESNgF-L9IrA7F7l8Ilt9ypkWv2SRiOuNImuCnQNngo7lDg6P99u2zoe7rWg0iIzmKilkv6KKKI_A"; //=> tungdt@83.com.vn
// $tokenTest = "1//01qvD-6luyaMaCgYIARAAGAESNwF-L9IrGzQGdhGCvLWXc1HV-RO6sk2RfosQZVce3F65IFsltfBNQHNVCZdxWU6CWsPq9NgvNeA"; //=> tungdtdev@gmail.com
if (empty($tokenTest)) {
    $val = getopt(null, ["refreshToken:"]);
    if (empty($val["refreshToken"])) {
        echo json_encode(["status" => false, "message" => "refresh token is required"]);
        exit();
    }
    $refreshToken = $val["refreshToken"];
} else {
    $refreshToken = $tokenTest;
}

require __DIR__ . "/googleads-php-lib/vendor/autoload.php";

use Google\Auth\OAuth2;
use Google\AdsApi\AdManager\AdManagerSessionBuilder;

//https://myaccount.google.com/u/0/permissions?pli=1
$oauth2 = new OAuth2([
    'authorizationUri' => 'https://accounts.google.com/o/oauth2/v2/auth',
    'tokenCredentialUri' => 'https://www.googleapis.com/oauth2/v4/token',
    'redirectUri' => "https://apps.valueimpression.com/gam/callback",
    //     'clientId' => '697051439959-ola1rso2vsjf3cqhvmnoh3lnmkd5mn1t.apps.googleusercontent.com',
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
    echo json_encode(["status" => false, "message" => $exc->getMessage()]);
    //     var_dump($exc->getMessage());
    //     var_dump($exc->getCode());
    //     var_dump($exc->getLine());
    exit();
}

try {
    $fileIni = __DIR__ . "/adsapi_php.ini";
    $session = (new AdManagerSessionBuilder())
        ->fromFile($fileIni)
        ->withOAuth2Credential($oauth2)
        ->build();
} catch (Exception $exc) {
    echo $exc->getTraceAsString();
    exit();
}
