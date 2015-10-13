<?php

function connect_AD()
  {
    $ldap_server = "ldaps://10.20.12.20";
    $ldap_user   = "CN=Administrator,CN=Users,DC=intra,DC=nanocloud,DC=com" ;
    $ldap_pass   = "password" ;

    $ad = ldap_connect($ldap_server) ;
    ldap_set_option($ad, LDAP_OPT_PROTOCOL_VERSION, 3) ;
    $bound = ldap_bind($ad, $ldap_user, $ldap_pass) ;

    return $ad ;
  }

    $ad = connect_AD();
    $cn = $argv[1];
    $password = $argv[2];

    $ldaprecord["CN"] = $cn;
    $ldaprecord["givenName"] = $cn;
    $ldaprecord["userPrincipalName"] = $cn;
    $ldaprecord["objectClass"] = "User";
    $ldaprecord["unicodePwd"] = mb_convert_encoding('"' . $password . '"', 'utf-16le');
    $ldaprecord["UserAccountControl"] = "512";


    $dn = "CN=$cn,OU=TalendUsers,DC=intra,DC=nxbay,DC=com";

    $result = ldap_add($ad, $dn, $ldaprecord);

    //$modifyUser["samAccountName"] = "";//$cn;
    //$result = ldap_modify($ad, $dn, $modifyUser); 

    $sr = ldap_search($ad,"OU=TalendUsers,DC=intra,DC=nxbay,DC=com","cn=$cn");
    $info = ldap_get_entries($ad,$sr);
    echo $info[0]["samaccountname"][0];

?>
