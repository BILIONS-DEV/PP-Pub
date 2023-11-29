<?php

include("config.php");

use Google\AdsApi\AdManager\Util\v202302\StatementBuilder;
use Google\AdsApi\AdManager\v202302\ServiceFactory;

$serviceFactory = new ServiceFactory();


$networkService = $serviceFactory->createNetworkService($session);
$networks = $networkService->getAllNetworks();
var_dump($networks);
die;


$userService = $serviceFactory->createUserService($session);

$roles = $userService->getAllRoles();

// Print out some information for each role.
foreach ($roles as $i => $role) {
    printf(
        "%d) Role with ID %d and name '%s' was found.%s",
        $i,
        $role->getId(),
        $role->getName(),
        PHP_EOL
    );
}

printf("Number of results found: %d%s", count($roles), PHP_EOL);
die;

try {
    // Get the current user.
    $user = $userService->getCurrentUser();
} catch (\Google\AdsApi\AdManager\v202302\ApiException $e) {
    echo json_encode(["status" => false, "message" => "Unable to get GAM manager information"]);
    exit();
}
if (empty($user)) {
    echo json_encode(["status" => false, "message" => "No accessible user of dfp found"]);
    return;
}

$_user = [
    "id" => $user->getId(),
    "name" => $user->getName(),
    "email" => $user->getEmail(),
    "role" => $user->getRoleName()
];

echo json_encode(["status" => true, "message" => "", "user" => $_user], JSON_NUMERIC_CHECK);
