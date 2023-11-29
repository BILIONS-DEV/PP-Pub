<?php

include("config.php");

use Google\AdsApi\AdManager\v202302\ServiceFactory;

$serviceFactory = new ServiceFactory();
$networkService = $serviceFactory->createNetworkService($session);
// Get all networks that you have access to with the current
// authentication credentials.
$networks = $networkService->getAllNetworks();
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
        "getChildPublishers" => $network->getChildPublishers()
    ];
    //     printf(
    //             "%d) Network with code %d and display name '%s' was found.%s", $i, $network->getNetworkCode(), $network->getDisplayName(), PHP_EOL
    //     );
}
echo json_encode(["status" => true, "message" => "", "networks" => $_networks], JSON_NUMERIC_CHECK);


//php network_backup.php --refreshToken=
//php network_backup.php --refreshToken=1//01jOjvnN36PAdCgYIARAAGAESNgF-L9IrVNVw89yVzIYHmw-JWqEt9F0wsWEWtORtI16H_Fg0lK07_O7BcF-QEHfTfvsk_XrzYQ