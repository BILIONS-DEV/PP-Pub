<?php

$val = getopt(null, ["key:", "value:"]);

if (empty($val["key"])) {
    echo json_encode(["status" => false, "message" => "key is required"]);
    return;
}
if ($val["key"]) {
    $key = $val["key"];
}

if (empty($val["value"])) {
    echo json_encode(["status" => false, "message" => "value is required"]);
    return;
}
if ($val["value"]) {
    $value = $val["value"];
}

$cpc = 0.0;


if (is_numeric($value)) {
    echo json_encode(["status" => true, "type" => "DECRYPT", "message" => "Success!", "data" => (float)$value]);
} else {
    $decoded_value = urldecode($value);
    //Convert Base64URL to Base64 by replacing "-" with "+" and "_" with "/"
    $b64 = strtr($decoded_value, '-_', '+/');
    $result = openssl_decrypt(base64_decode($b64), "AES-128-ECB", hex2bin($key), OPENSSL_PKCS1_PADDING);
    $validation_regex = "/^(\d)+((\.)\d+)?\_/";
    // Check that the decrypted value is valid
    if (preg_match($validation_regex, $result) == 0) {
        echo json_encode(["status" => false, "type" => "DECRYPT", "message" => 'Decrypted value is not valid!']);
        return;
    }
    list($cpc, $b) = explode("_", $result);
    echo json_encode(["status" => true, "type" => "DECRYPT", "message" => "Success!", "data" => (float)$cpc]);
}
?>
