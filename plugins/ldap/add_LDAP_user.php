<?php


include './connection.php';

$ldap_connection = connect_AD();

// Our DN
$ldap_base_dn = 'OU=NanocloudUsers,DC=intra,DC=nanocloud,DC=com';

// This filter will get all the users with disabled account
$search_filter = '(&(objectClass=User)(userAccountControl:1.2.840.113556.1.4.803:=2))';

// Enabled accounts
// '(&(objectClass=User)(!userAccountControl:1.2.840.113556.1.4.803:=2))'

$result = ldap_search($ldap_connection, $ldap_base_dn, $search_filter);

$count_disabled_account = ldap_count_entries($ldap_connection, $result);


if ($count_disabled_account) {

  $disabled_accounts = ldap_get_entries($ldap_connection, $result);

  $dn=$disabled_accounts[0]["dn"];

  $sam_account_name = $disabled_accounts[0]["samaccountname"][0];

  $ldaprecord["mail"] = $user_email;
  $ldaprecord["givenName"] = $user_email;
  $ldaprecord["userPrincipalName"] = $user_email;
  $ldaprecord["objectClass"] = "User";
  $ldaprecord["unicodePwd"] = mb_convert_encoding('"' . $password . '"', 'utf-16le');
  $ldaprecord["UserAccountControl"] = "512";

  // Update account
  $r = ldap_modify($ldap_connection, $dn, $ldaprecord);
  if ($r == FALSE) {
    fwrite(STDERR, "An error occurred during LDAP account update.\n");
    exit(1);
  }
}
else {

  // This filter will get all the users
  $search_filter = '(&(objectCategory=person)(samaccountname=*))';
  // Query the LDAP server
  $result1 = ldap_search($ldap_connection, $ldap_base_dn, $search_filter);
  $number_of_users = ldap_count_entries($ldap_connection, $result1);
  $cn = "demo" . sprintf('%04d', ++$number_of_users);

  $ldaprecord["CN"] = $cn;
  $ldaprecord["mail"] = $user_email;
  $ldaprecord["givenName"] = $cn;
  $ldaprecord["userPrincipalName"] = $cn;
  $ldaprecord["objectClass"] = "User";
  $ldaprecord["unicodePwd"] = mb_convert_encoding('"' . $password . '"', 'utf-16le');
  $ldaprecord["UserAccountControl"] = "512";


  $dn = "CN=$cn,OU=NanocloudUsers,DC=intra,DC=nanocloud,DC=com";

  // Insert new account
  $r = ldap_add($ldap_connection, $dn, $ldaprecord);
  if ($r == FALSE) {
    fwrite(STDERR, "An error occurred.\n");
    exit(1);
  }

  $sr = ldap_search($ldap_connection,"OU=NanocloudUsers,DC=intra,DC=nanocloud,DC=com","cn=$cn");
  $info = ldap_get_entries($ldap_connection,$sr);
  $sam_account_name =  $info[0]["samaccountname"][0];
}

ldap_unbind($ldap_connection) or die('Unable to close LDAP connection');

echo $sam_account_name;

?>
