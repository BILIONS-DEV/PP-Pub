<?php

include("config.php");

use Google\AdsApi\AdManager\v202302\ServiceFactory;

$serviceFactory = new ServiceFactory();

// GET USER
$userService = $serviceFactory->createUserService($session);
try {
    // Get the current user.
    $user = $userService->getCurrentUser();
    $_user = [
        "id" => $user->getId(),
        "name" => $user->getName(),
        "email" => $user->getEmail(),
        "role" => $user->getRoleName()
    ];
} catch (\Google\AdsApi\AdManager\v202302\ApiException $e) {
    //     echo json_encode(["status" => false, "message" => "Unable to get GAM manager information" .": ". $e->getMessage()]);
    $_user = [
        "id" => 0,
        "name" => "",
        "email" => "",
        "role" => ""
    ];
}

// GET NETWORK
$networkService = $serviceFactory->createNetworkService($session);
try {
    // Get all networks that you have access to with the current
    $networks = $networkService->getAllNetworks();
} catch (\Google\AdsApi\AdManager\v202302\ApiException $e) {
    echo json_encode(["status" => false, "message" => "Error get networks"]);
    exit();
}
if (empty($networks)) {
    echo json_encode(["status" => false, "message" => "No accessible networks found"]);
    return;
}
// Print out some information for each network.
$_networks = [];
foreach ($networks as $i => $network) {
    $_networks[] = [
        "id" => $network->getNetworkCode(),
        "name" => $network->getDisplayName(),
        "currencyCode" => $network->getCurrencyCode(),
        "getChildPublishers" => $network->getChildPublishers(),
        "timeZone" => $network->getTimeZone()
    ];
}
echo json_encode(["status" => true, "message" => "", "user" => $_user, "networks" => $_networks], JSON_NUMERIC_CHECK);
exit();
//php network_backup.php --refreshToken=
//php network_backup.php --refreshToken=1//01jOjvnN36PAdCgYIARAAGAESNgF-L9IrVNVw89yVzIYHmw-JWqEt9F0wsWEWtORtI16H_Fg0lK07_O7BcF-QEHfTfvsk_XrzYQ