import * as React from 'react';
import { styled, createTheme, ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import MuiDrawer from '@mui/material/Drawer';
import Box from '@mui/material/Box';
import MuiAppBar, { AppBarProps as MuiAppBarProps } from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import List from '@mui/material/List';
import Typography from '@mui/material/Typography';
import { mainListItems } from './menu';
import { Grid, MenuItem, Select } from '@mui/material';
import { AppBarContext } from './contexts/AppBarContext';
import { Link } from 'react-router-dom';

import blackdaggerImage from './assets/images/blackdagger.png';

const drawerWidthClosed = 64;
const drawerWidth = 240;

interface AppBarProps extends MuiAppBarProps {
  open?: boolean;
}

const AppBar = styled(MuiAppBar, {
  shouldForwardProp: (prop) => prop !== 'open',
})<AppBarProps>(({ theme, open }) => ({
  zIndex: theme.zIndex.drawer - 1,
  transition: theme.transitions.create(['width', 'margin', 'border'], {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen,
  }),
  width: '100%',
  ...(open && {
    transition: theme.transitions.create(['width', 'margin', 'border'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
  }),
}));

const Drawer = styled(MuiDrawer, {
  shouldForwardProp: (prop) => prop !== 'open',
})(({ theme, open }) => ({
  '& .MuiDrawer-paper': {
    position: 'relative',
    whiteSpace: 'nowrap',
    width: drawerWidth,
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
    boxSizing: 'border-box',
    ...(!open && {
      overflowX: 'hidden',
      transition: theme.transitions.create('width', {
        easing: theme.transitions.easing.sharp,
        duration: theme.transitions.duration.leavingScreen,
      }),
      width: drawerWidthClosed,
      [theme.breakpoints.up('sm')]: {
        width: theme.spacing(9),
      },
    }),
  },
}));

const mdTheme = createTheme({
  palette: {
    mode: 'dark', // Set the theme mode to dark
    background: {
      default: '#121212', // Dark background for most surfaces
      paper: '#1e1e1e', // Dark background for components like Drawer and AppBar
    },
    text: {
      primary: '#FFFEFE', // Light text color for readability
      secondary: '#FFFEFE', // Slightly dimmer text color for less emphasis
    },
  },
  typography: {
    fontFamily: 'Inter',
  },
});

type DashboardContentProps = {
  title: string;
  navbarColor: string;
  version: string;
  children?: React.ReactElement | React.ReactElement[];
};

function Content({ title, navbarColor, children }: DashboardContentProps) {
  const [scrolled, setScrolled] = React.useState(false);
  const containerRef = React.useRef<HTMLDivElement>(null);
  const gradientColor = navbarColor || '#171617';

  return (
    <ThemeProvider theme={mdTheme}>
      <Box sx={{ display: 'flex', flexDirection: 'row', width: '100vw' }}>
        <CssBaseline />
        <Drawer variant="permanent" open={false}>
          <Box
            sx={{
              background: `linear-gradient(0deg, ${mdTheme.palette.background.default} 0%, ${gradientColor} 70%, ${gradientColor} 100%)`,
              height: '100%',
            }}
          >
            <List component="nav" sx={{ pl: '6px' }}>
              {mainListItems}
            </List>
          </Box>
        </Drawer>
        <Box
          component="main"
          sx={{
            display: 'flex',
            flexDirection: 'column',
            backgroundColor: 'white',
            height: '100vh',
            width: '100%',
            maxWidth: '100%',
            overflow: 'auto',
          }}
        >
          <AppBar
            open={false}
            elevation={0}
            sx={{
              width: '100%',
              backgroundColor: '#171617',
              borderBottom: scrolled ? 1 : 0,
              borderColor: 'grey.300',
              pr: 3,
              position: 'relative',
              display: 'block',
            }}
          >
            <Toolbar
              sx={{
                width: '100%',
                display: 'flex',
                direction: 'row',
                justifyContent: 'space-between',
                alignItems: 'center',
                flex: 1,
              }}
            >
              <AppBarContext.Consumer>
                {(context) => (
                  <NavBarTitleText visible={scrolled}>
                    {context.title}
                  </NavBarTitleText>
                )}
              </AppBarContext.Consumer>

              <Box sx={{ display: 'flex', alignItems: 'center' }}>
                <Link to="/dashboard">
                  <NavBarTitleText>{title || 'Blackdagger'}</NavBarTitleText>
                </Link>
                <Link to="/dashboard">
                  <img src={blackdaggerImage} alt="Blackdagger" />
                </Link>

                <AppBarContext.Consumer>
                  {(context) =>
                    context.remoteNodes && context.remoteNodes.length > 0 ? (
                      <Select
                        sx={{
                          backgroundColor: 'white',
                          color: 'black',
                          borderRadius: '5px',
                          border: '1px solid #ccc',
                          marginLeft: '10px',
                          height: '30px',
                          width: '150px',
                          marginBottom: '5px',
                        }}
                        value={context.selectedRemoteNode}
                        onChange={(e) => {
                          context.selectRemoteNode(e.target.value);
                        }}
                      >
                        {context.remoteNodes.map((node) => (
                          <MenuItem key={node} value={node}>
                            {node}
                          </MenuItem>
                        ))}
                      </Select>
                    ) : null
                  }
                </AppBarContext.Consumer>
              </Box>
            </Toolbar>
          </AppBar>

          <Grid
            container
            ref={containerRef}
            sx={{
              flex: 1,
              pb: 4,
              overflow: 'auto',
              backgroundColor: '#171617',
            }}
            onScroll={() => {
              const curr = containerRef.current;
              if (curr) {
                setScrolled(curr.scrollTop > 54);
              }
            }}
          >
            {children}
          </Grid>
        </Box>
      </Box>
    </ThemeProvider>
  );
}

type NavBarTitleTextProps = {
  children: string;
  visible?: boolean;
};

const NavBarTitleText = ({
  children,
  visible = true,
}: NavBarTitleTextProps) => (
  <Typography
    component="h1"
    variant="h6"
    gutterBottom
    sx={{
      fontWeight: '800',
      color: '#FFFEFE',
      opacity: visible ? 1 : 0,
      transition: 'opacity 0.2s',
    }}
  >
    {children}
  </Typography>
);

type DashboardProps = DashboardContentProps;

export default function Layout({ children, ...props }: DashboardProps) {
  return <Content {...props}>{children}</Content>;
}
