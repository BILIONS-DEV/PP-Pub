<?php

include("config.php");

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


echo json_encode(["status" => true, "type" => "GET_ALL_ACCOUNTS", "message" => "Success!", "data" => $accounts]);