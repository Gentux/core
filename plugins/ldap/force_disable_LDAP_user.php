<?php

// Configuration
$ldap_server = "ldaps://10.20.12.20";
$ldap_user   = "CN=Administrator,CN=Users,DC=intra,DC=nanocloud,DC=com";
$ldap_pass   = "password";

// Command line parameters
$email = $argv[1];

// Connection
$ldap_connection = ldap_connect($ldap_server) or die('Unable to connect to LDAP server');

// We have to set this option for the version of Active Directory we are using.
ldap_set_option($ldap_connection, LDAP_OPT_PROTOCOL_VERSION, 3) or die('Unable to set LDAP protocol version');
ldap_set_option($ldap_connection, LDAP_OPT_REFERRALS, 0); // We need this for doing an LDAP search.    

// Binding
ldap_bind($ldap_connection, $ldap_user, $ldap_pass) or die('Unable to bind to LDAP server');

// Our DN
$ldap_base_dn = 'OU=NanocloudUsers,DC=intra,DC=nanocloud,DC=com';

// This filter will get the user
$search_filter = '(&(objectCategory=person)(mail=' . $email . ')(!(userAccountControl:1.2.840.113556.1.4.803:=2)))';

$result = ldap_search($ldap_connection, $ldap_base_dn, $search_filter);

$count_accounts = ldap_count_entries($ldap_connection, $result);

if ($count_accounts == 1) {

  $account = ldap_get_entries($ldap_connection, $result);

  $dn=$account[0]["dn"];

  $ldaprecord["objectClass"] = "User";
  $ldaprecord["UserAccountControl"] = "514";

  // Update account
  $r = ldap_modify($ldap_connection, $dn, $ldaprecord);

  if ($r == FALSE) {
    fwrite(STDERR, "An error occurred during LDAP modification\n");
    exit(1);
  }
}
else {
  $count_accounts = 0;
  fwrite(STDERR, "An error occurred. $email account not available\n");
}

ldap_unbind($ldap_connection) or die('Unable to close LDAP connection');

echo $count_accounts;
?>
