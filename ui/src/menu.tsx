import * as React from 'react';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import { Link } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
  faChartGantt,
  faMagnifyingGlass,
  faTableList,
  faTerminal,
  faBook,
} from '@fortawesome/free-solid-svg-icons';
import { IconProp } from '@fortawesome/fontawesome-svg-core';
import { Typography } from '@mui/material';
import Tooltip from '@mui/material/Tooltip';

function Icon({ icon }: { icon: IconProp }) {
  return (
    <span
      style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        marginLeft: 2,
      }}
    >
      <FontAwesomeIcon
        style={{ height: 20, width: 20, color: '#FFFEFE' }}
        icon={icon}
      ></FontAwesomeIcon>
    </span>
  );
}

export const mainListItems = (
  <React.Fragment>
    <Link to="/dashboard">
      <ListItem text="Dashboard" icon={<Icon icon={faChartGantt} />} />
    </Link>
    <Link to="/dags">
      <ListItem text="DAGs" icon={<Icon icon={faTableList} />} />
    </Link>
    <Link to="/search">
      <ListItem text="Search" icon={<Icon icon={faMagnifyingGlass} />} />
    </Link>
    <ListItem
      text="Terminal"
      icon={<Icon icon={faTerminal} />}
      to="http://:8090"
      external
    />
    <ListItemDoc
      text="Documentation"
      icon={<Icon icon={faBook} />}
      to="https://blackdagger.readthedocs.io/en/latest/"
      external
    />
  </React.Fragment>
);

type ListItemProps = {
  icon: React.ReactNode;
  text: string;
  to?: string;
  external?: boolean;
};

function ListItem({ icon, text, to, external }: ListItemProps) {
  let listItemProps = {};

  if (external) {
    const protocol = window.location.protocol;
    const hostname = window.location.hostname;
    const port = '8090';
    const href = `${protocol}//${hostname}:${port}`;
    listItemProps = {
      component: 'a',
      href: href,
      target: '_blank',
      rel: 'noopener noreferrer',
    };
  } else {
    listItemProps = { component: 'a', href: to };
  }

  const content = (
    <ListItemButton component="a" {...listItemProps}>
      <ListItemIcon sx={{ color: 'black' }}>{icon}</ListItemIcon>
      <ListItemText
        primary={
          <Typography
            sx={{
              color: 'white',
              fontWeight: '400',
            }}
          >
            {text}
          </Typography>
        }
      />
    </ListItemButton>
  );

  // If the list item is for the Terminal, wrap it in a Tooltip
  if (text === 'Terminal') {
    return (
      <Tooltip
        title="If you want to access the terminal interface please make sure `default-gotty-service` dag is running"
        arrow
      >
        {content}
      </Tooltip>
    );
  }

  return content;
}

function ListItemDoc({ icon, text, to, external }: ListItemProps) {
  let listItemProps = {};

  if (external) {
    // Directly use the 'to' prop for external links
    listItemProps = {
      component: 'a',
      href: to,
      target: '_blank',
      rel: 'noopener noreferrer',
    };
  } else {
    // For internal routing, adjust as needed for your routing library
    listItemProps = { component: 'a', href: to };
  }
  return (
    <ListItemButton component="a" {...listItemProps}>
      <ListItemIcon sx={{ color: 'black' }}>{icon}</ListItemIcon>
      <ListItemText
        primary={
          <Typography
            sx={{
              color: 'white',
              fontWeight: '400',
            }}
          >
            {text}
          </Typography>
        }
      />
    </ListItemButton>
  );
}
