<?php

// Configuration
$ldap_server = "ldaps://10.20.12.20";
$ldap_user   = "CN=Administrator,CN=Users,DC=intra,DC=nanocloud,DC=com";
$ldap_pass   = "password";

// Connection
$ldap_connection = ldap_connect($ldap_server) or die('Unable to connect to LDAP server');

// We have to set this option for the version of Active Directory we are using.
ldap_set_option($ldap_connection, LDAP_OPT_PROTOCOL_VERSION, 3) or die('Unable to set LDAP protocol version');
ldap_set_option($ldap_connection, LDAP_OPT_REFERRALS, 0); // We need this for doing an LDAP search.    

// Binding
ldap_bind($ldap_connection, $ldap_user, $ldap_pass) or die('Unable to bind to LDAP server');

// Our DN
$ldap_base_dn = 'OU=TalendUsers,DC=intra,DC=nanocloud,DC=com';

// This filter will get all the users with disabled account
$search_filter = '(&(objectClass=User)(userAccountControl:1.2.840.113556.1.4.803:=2))';

// Enabled accounts
// '(&(objectClass=User)(!userAccountControl:1.2.840.113556.1.4.803:=2))'

// Query the LDAP server
$result = ldap_search($ldap_connection, $ldap_base_dn, $search_filter);

echo "Number of entries returned is " . ldap_count_entries($ldap_connection, $result) . "\n";

echo "Getting entries ...\n";
$info = ldap_get_entries($ldap_connection, $result);
echo "Data for " . $info["count"] . " items returned:\n";

for ($i=0; $i<$info["count"]; $i++) {
  echo "--------------------------------------------\n";
  echo "dn is: " . $info[$i]["dn"] . "\n";
  echo "first cn entry is: " . $info[$i]["cn"][0] . "\n";
  echo "Samaccountname is: " . $info[$i]["samaccountname"][0] . "\n";
  echo "UserAccountControl is: " . $info[$i]["useraccountcontrol"][0] . "\n";
  $ac = $info[$i]["useraccountcontrol"][0];
  if (($ac & 2)==2) $status="Disabled"; else $status="Enabled";
  echo "User is " . $status . "\n";
}

// Leaving the LDAP server
ldap_unbind($ldap_connection) or die('Unable to close LDAP connection');

?>
